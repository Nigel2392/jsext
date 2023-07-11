package jsc

import (
	"encoding/base64"
	"reflect"
	"strings"
	"syscall/js"

	"github.com/Nigel2392/jsext/v2"
	"github.com/Nigel2392/jsext/v2/jse"
)

// TINYGO is used to check if we are using tinygo as a compiler.
var TINYGO bool

// Package JSC implements a way to convert javascript objects to go objects, and vice versa.
//
// This package is used to communicate between the frontend and backend.
//
// If the js tag is not specified on struct fields, it will default to using json.
//
// This file uses reflection, thus might cause unforseen errors with TinyGo.
//
// When using tinygo as a compiler, be sure to set jsc.TINYGO = true.
// This will try to make sure no unforseen panics happen during runtime.

type Marshaller interface {
	MarshalJS() js.Value
}

func ValueOf(f any) js.Value {
	if f == nil {
		return js.Null()
	}
	switch val := f.(type) {
	case int, int64, int32, int16, int8,
		float64, float32,
		uint, uint64, uint32, uint16, uint8, uintptr,
		string, bool:
		// []any, map[string]any: // Removed so we can call jss.ValueOf on a slice or map.

		return js.ValueOf(val)
	case js.Value, js.Func:
		return js.ValueOf(val)
	case jsext.Value:
		return val.Value()
	case *jse.Element:
		return val.JSValue()
	case jsext.Element:
		return val.JSValue()
	case jsext.Event:
		return val.JSValue()
	case jsext.Import:
		return val.JSValue()
	case jsext.Promise:
		return val.JSValue()
	case []byte:
		var enc = base64.StdEncoding.EncodeToString(val)
		return js.ValueOf(enc)
	case func():
		if val == nil {
			return js.Null()
		}
		return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			val()
			return nil
		}).Value
	case func(this js.Value, args []js.Value) interface{}:
		if val == nil {
			return js.Null()
		}
		return js.FuncOf(val).Value
	case Marshaller:
		return val.MarshalJS()
	}
	var valueOf = reflect.ValueOf(f)
	var kind = valueOf.Kind()
	switch kind {
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return js.ValueOf(valueOf.Int())
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uintptr:
		return js.ValueOf(valueOf.Uint())
	case reflect.Float64, reflect.Float32:
		return js.ValueOf(valueOf.Float())
	case reflect.String:
		return js.ValueOf(valueOf.String())
	case reflect.Bool:
		return js.ValueOf(valueOf.Bool())
	case reflect.Slice, reflect.Array:
		if !valueOf.IsValid() {
			return js.Null()
		}

		// Check if bytes
		if valueOf.Type().Elem().Kind() == reflect.Uint8 {
			var enc = base64.StdEncoding.EncodeToString(valueOf.Bytes())
			return js.ValueOf(enc)
		}

		var length = valueOf.Len()
		var array = js.Global().Get("Array").New(length)
		for i := 0; i < length; i++ {
			var index = valueOf.Index(i)
			if index.Kind() == reflect.Ptr {
				index = index.Elem()
			}
			array.SetIndex(i, ValueOf(index.Interface()))
		}
		return js.ValueOf(array)
	case reflect.Map:
		var keys = valueOf.MapKeys()
		var object = js.Global().Get("Object").New()
		for _, key := range keys {
			object.Set(key.String(), ValueOf(valueOf.MapIndex(key).Interface()))
		}
		return object
	case reflect.Struct:
		if !valueOf.CanInterface() {
			return js.Null()
		}
		if !TINYGO {
			if valueOf.Type().ConvertibleTo(reflect.TypeOf(js.Value{})) {
				var jsValue js.Value
				reflect.ValueOf(&jsValue).Elem().Set(valueOf)
				return jsValue
			}
		}
		var object = js.Global().Get("Object").New()
		var typeOf = valueOf.Type()
		var numField = valueOf.NumField()
		for i := 0; i < numField; i++ {
			var field = typeOf.Field(i)
			var tag, omitEmpty, ok = getStructTag(field, "js", "json", "jsc")
			if !ok {
				continue
			}
			var valField = valueOf.Field(i)
			if valField.Kind() == reflect.Ptr {
				valField = valField.Elem()
			}

			if omitEmpty && valField.IsZero() {
				continue
			}

			if !valField.CanInterface() {
				panic("ValueOf: cannot interface " + valField.String())
			}

			object.Set(tag, ValueOf(valField.Interface()))
		}
		return object
	case reflect.Ptr:
		if valueOf.IsNil() {
			return js.Null()
		}
		return ValueOf(valueOf.Elem().Interface())
	default:
		panic("ValueOf: unsupported type " + kind.String())
	}
}

func getStructTag(field reflect.StructField, tags ...string) (name string, omitEmpty bool, ok bool) {
	for _, tag := range tags {
		var value = field.Tag.Get(tag)
		if value != "" {
			name = value
			break
		}
	}
	if name == "" {
		name = field.Name
	}
	if name == "-" {
		return "", false, false
	}

	ok = true
	omitEmpty = false

	if strings.Index(name, ",") != -1 {
		var parts = strings.Split(name, ",")
		name = parts[0]
		omitEmpty = strings.ToLower(parts[1]) == "omitempty"
	}

	return name, omitEmpty, ok
}

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrUndefined Error = "src is null or undefined"
	ErrNil       Error = "dst is nil"
	ErrNotObject Error = "src is not an object"
	ErrNotPtr    Error = "dst is not a pointer"
	ErrNotValid  Error = "dst is not a pointer to a struct, map, or slice"
	ErrNotStruct Error = "dst is not a pointer to a struct"
	ErrCannotSet Error = "cannot set dst field"
)

type Unmarshaller interface {
	UnmarshalJS(js.Value) error
}

func Scan(src js.Value, dst interface{}) error {
	if src.IsNull() || src.IsUndefined() {
		return ErrUndefined
	}

	if dst == nil {
		return ErrNil
	}

	if src.Type() != js.TypeObject {
		return ErrNotObject
	}

	var (
		dstVal = reflect.ValueOf(dst)
		dstTyp = dstVal.Type()
	)

	if dstTyp.Kind() != reflect.Ptr {
		return ErrNotPtr
	}

	dstVal = dstVal.Elem()

	return scanValue(src, dstVal)
}

func scanStruct(src js.Value, dstVal reflect.Value, dstTyp reflect.Type) error {
	if dstTyp.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	var numField = dstTyp.NumField()
	for i := 0; i < numField; i++ {
		var dstField = dstTyp.Field(i)
		var dstTag, _, dstOk = getStructTag(dstField, "js", "json", "jsc")
		if !dstOk {
			continue
		}
		var dstValField = dstVal.Field(i)
		if dstValField.Kind() == reflect.Ptr {
			dstValField = dstValField.Elem()
		}
		if !dstValField.CanSet() {
			return ErrCannotSet
		}
		var srcVal = src.Get(dstTag)
		if srcVal.IsUndefined() || srcVal.IsNull() {
			continue
		}
		var err = scanValue(srcVal, dstValField)
		if err != nil {
			return err
		}
	}
	return nil
}

var unmarshallerType = reflect.TypeOf((*Unmarshaller)(nil)).Elem()

func scanValue(srcVal js.Value, dstVal reflect.Value) error {
	if dstVal.Kind() == reflect.Ptr {
		dstVal = dstVal.Elem()
	}

	if dstVal.CanAddr() && dstVal.Addr().Type().Implements(unmarshallerType) {
		var unmarshaller = dstVal.Addr().Interface().(Unmarshaller)
		return unmarshaller.UnmarshalJS(srcVal)
	}

	switch dstVal.Kind() {
	case reflect.String:
		dstVal.SetString(srcVal.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dstVal.SetInt(int64(srcVal.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dstVal.SetUint(uint64(srcVal.Int()))
	case reflect.Float32, reflect.Float64:
		dstVal.SetFloat(srcVal.Float())
	case reflect.Bool:
		dstVal.SetBool(srcVal.Bool())
	case reflect.Slice:
		if dstVal.Type().Elem().Kind() == reflect.Uint8 {
			var bytes, err = base64.StdEncoding.DecodeString(srcVal.String())
			if err == nil {
				dstVal.SetBytes(bytes)
			} else {
				var b = make([]byte, srcVal.Length())
				js.CopyBytesToGo(b, srcVal)
			}
			return nil
		}
		var err = scanSlice(srcVal, dstVal)
		if err != nil {
			return err
		}
	case reflect.Struct:
		var err = scanStruct(srcVal, dstVal, dstVal.Type())
		if err != nil {
			return err
		}
	case reflect.Map:
		var err = scanMap(srcVal, dstVal)
		if err != nil {
			return err
		}
	case reflect.Interface:
		dstVal.Set(reflect.ValueOf(guessType(srcVal)))
	case reflect.Ptr:
		var err = scanValue(srcVal, dstVal.Elem())
		if err != nil {
			return err
		}
	}
	return nil
}

func guessType(srcVal js.Value) interface{} {
	var i interface{}
	switch srcVal.Type() {
	case js.TypeBoolean:
		i = srcVal.Bool()
	case js.TypeString:
		i = srcVal.String()
	case js.TypeNumber:
		i = srcVal.Float()
	case js.TypeObject:
		if srcVal.InstanceOf(js.Global().Get("Array")) {
			var s = make([]interface{}, srcVal.Length())
			for i := 0; i < srcVal.Length(); i++ {
				s[i] = guessType(srcVal.Index(i))
			}
			i = s
			break
		}
		var m = make(map[string]interface{})
		var valueOf = reflect.ValueOf(&m)
		var err = scanMap(srcVal, valueOf)
		if err != nil {
			return err
		}
		i = m
	case js.TypeFunction:
		i = func(args ...any) js.Value {
			var s = make([]interface{}, len(args))
			for i, arg := range args {
				s[i] = ValueOf(arg)
			}
			return srcVal.Invoke(s...)
		}
	}
	return i
}

func scanMap(srcVal js.Value, dstVal reflect.Value) error {
	if dstVal.IsNil() {
		dstVal.Set(reflect.MakeMap(dstVal.Type()))
	}
	if dstVal.Kind() == reflect.Ptr {
		dstVal = dstVal.Elem()
	}
	var keys = js.Global().Get("Object").Call("keys", srcVal)
	var numKeys = keys.Length()
	for i := 0; i < numKeys; i++ {
		var srcKey = keys.Index(i)
		var srcKeyValue = srcVal.Get(srcKey.String())
		var dstKey = reflect.New(dstVal.Type().Key())
		var err = scanValue(srcKey, dstKey)
		if err != nil {
			return err
		}

		var scanInto = dstVal.MapIndex(dstKey.Elem())
		if scanInto.IsValid() {
			err = scanValue(srcKeyValue, scanInto)
			if err != nil {
				return err
			}
			continue
		}

		var typeOfDstValElem = dstVal.Type().Elem()
		if typeOfDstValElem.Kind() == reflect.Ptr {
			typeOfDstValElem = typeOfDstValElem.Elem()
		}
		var dstKeyValue = reflect.New(typeOfDstValElem)

		err = scanValue(srcKeyValue, dstKeyValue)
		if err != nil {
			return err
		}

		dstVal.SetMapIndex(dstKey.Elem(), dstKeyValue.Elem())
	}
	return nil
}

func scanSlice(srcVal js.Value, dstVal reflect.Value) error {
	var srcLen = srcVal.Length()
	if srcLen == 0 {
		return nil
	}
	if dstVal.IsNil() {
		// makeslice is implemented in tinygo! :)
		dstVal.Set(reflect.MakeSlice(dstVal.Type(), srcLen, srcLen))
	}
	for i := 0; i < srcLen; i++ {
		var srcElem = srcVal.Index(i)
		var dstElem = reflect.New(dstVal.Type().Elem())
		var err = scanValue(srcElem, dstElem)
		if err != nil {
			return err
		}
		// slice index out of range
		if i >= dstVal.Len() {
			dstVal.Set(reflect.Append(dstVal, dstElem.Elem()))
		} else {
			dstVal.Index(i).Set(dstElem.Elem())
		}
	}
	return nil
}

// makeslice is implemented in tinygo! :)
//	func tinySetNilSlice(dstVal reflect.Value, elemKind reflect.Kind, srcLen int) error {
//		switch elemKind {
//		case reflect.Interface:
//			dstVal.Set(reflect.ValueOf(make([]interface{}, srcLen)))
//		case reflect.String:
//			dstVal.Set(reflect.ValueOf(make([]string, srcLen)))
//		case reflect.Int:
//			dstVal.Set(reflect.ValueOf(make([]int, srcLen)))
//		case reflect.Int8:
//			dstVal.Set(reflect.ValueOf(make([]int8, srcLen)))
//		case reflect.Int16:
//			dstVal.Set(reflect.ValueOf(make([]int16, srcLen)))
//		case reflect.Int32:
//			dstVal.Set(reflect.ValueOf(make([]int32, srcLen)))
//		case reflect.Int64:
//			dstVal.Set(reflect.ValueOf(make([]int64, srcLen)))
//		case reflect.Uint:
//			dstVal.Set(reflect.ValueOf(make([]uint, srcLen)))
//		case reflect.Uint8:
//			dstVal.Set(reflect.ValueOf(make([]uint8, srcLen)))
//		case reflect.Uint16:
//			dstVal.Set(reflect.ValueOf(make([]uint16, srcLen)))
//		case reflect.Uint32:
//			dstVal.Set(reflect.ValueOf(make([]uint32, srcLen)))
//		case reflect.Uint64:
//			dstVal.Set(reflect.ValueOf(make([]uint64, srcLen)))
//		case reflect.Float32:
//			dstVal.Set(reflect.ValueOf(make([]float32, srcLen)))
//		case reflect.Float64:
//			dstVal.Set(reflect.ValueOf(make([]float64, srcLen)))
//		case reflect.Bool:
//			dstVal.Set(reflect.ValueOf(make([]bool, srcLen)))
//		case reflect.Complex64:
//			dstVal.Set(reflect.ValueOf(make([]complex64, srcLen)))
//		case reflect.Complex128:
//			dstVal.Set(reflect.ValueOf(make([]complex128, srcLen)))
//		case reflect.Uintptr:
//			dstVal.Set(reflect.ValueOf(make([]uintptr, srcLen)))
//		case reflect.UnsafePointer:
//			dstVal.Set(reflect.ValueOf(make([]unsafe.Pointer, srcLen)))
//		default:
//			return fmt.Errorf("tinygo: unsupported slice element type: %v", elemKind)
//		}
//		return nil
//	}
//
