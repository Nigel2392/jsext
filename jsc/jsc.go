package jsc

import (
	"encoding/base64"
	"reflect"
	"strings"
	"syscall/js"

	"github.com/Nigel2392/jsext/v2"
	"github.com/Nigel2392/jsext/v2/console"
	"github.com/Nigel2392/jsext/v2/errs"
	"github.com/Nigel2392/jsext/v2/jse"
)

// TINYGO is used to check if we are using tinygo as a compiler.
var TINYGO bool

// Wether to encode and decode bytes using base64.
var BASE64 = true

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

type ErrorMarshaller interface {
	MarshalJS() (js.Value, error)
}

/*

	I AM UNSURE OF THE USE FOR THE FOLLOWING FUNCTION, LEAVING IT OUT.

	type ref uint64

	// Value represents a JavaScript value. The zero value is the JavaScript value "undefined".
	// Values can be checked for equality with the Equal method.
	type jsValue struct {
		_     [0]func() // uncomparable; to make == not compile
		ref   ref       // identifies a JavaScript value, see ref type
		gcPtr *ref      // used to trigger the finalizer when the Value is not referenced any more
	}

	// Swap swaps the underlying value of two js.Value objects.
	func Swap[T js.Value | jsext.Value | jsext.Element | jsext.Event | jsext.Import | jse.Element](a, b *T) {
		var (
			aTyp, bTyp     = (*jsValue)(unsafe.Pointer(a)), (*jsValue)(unsafe.Pointer(b))
			aRef, bRef     = aTyp.ref, bTyp.ref
			aGCPtr, bGCPtr = aTyp.gcPtr, bTyp.gcPtr
		)

		aTyp.ref, bTyp.ref = bRef, aRef
		aTyp.gcPtr, bTyp.gcPtr = bGCPtr, aGCPtr

		runtime.KeepAlive(a)
		runtime.KeepAlive(b)
	}
*/

// ValuesOf will return the js.Value of the given values.
// It will return an error if interface{} of the value are not supported.
func ValuesOf(f ...interface{}) ([]js.Value, error) {
	var (
		v      js.Value
		err    error
		values = make([]js.Value, len(f))
	)
	for i := range f {
		v, err = ValueOf(f[i])
		if err != nil {
			return nil, err
		}
		values[i] = v
	}
	return values, nil
}

// ValuesOfInterface will return the js.Value of the given values in an array of interface{}.
// It will return an error if interface{} of the value are not supported.
//
// This is a helper function, since js.(Value).Call or set will not take a slice of js.Value.
func ValuesOfInterface(f ...interface{}) ([]interface{}, error) {
	var (
		v      js.Value
		err    error
		values = make([]interface{}, len(f))
	)
	for i := range f {
		v, err = ValueOf(f[i])
		if err != nil {
			return nil, err
		}
		values[i] = v
	}
	return values, nil
}

// MustValuesOf will return the js.Value of the given values.
// This can be useful for inlining known-safe types.
func MustValuesOf(f ...interface{}) []js.Value {
	return mustValuesOf(ValuesOf, f...)
}

// MustValuesOfInterface will return the js.Value of the given values in an array of interface{}.
// This can be useful for inlining known-safe types.
//
// This is a helper function, since js.(Value).Call or set will not take a slice of js.Value.
func MustValuesOfInterface(f ...interface{}) []interface{} {
	return mustValuesOf(ValuesOfInterface, f...)
}

func mustValuesOf[T interface{} | js.Value](f func(args ...interface{}) ([]T, error), args ...interface{}) []T {
	var v, err = f(args...)
	if err != nil {
		panic(err)
	}
	return v
}

// ValueOf will return the js.Value of the given value.
// It will return an error if the value is not supported.
func ValueOf(f interface{}) (js.Value, error) {
	if f == nil {
		return js.Null(), nil
	}
	switch val := f.(type) {
	case js.Value, js.Func:
		return js.ValueOf(val), nil
	case jsext.Value:
		return val.Value(), nil
	case *jse.Element:
		return val.JSValue(), nil
	case jsext.Element:
		return val.JSValue(), nil
	case jsext.Event:
		return val.JSValue(), nil
	case jsext.Import:
		return val.JSValue(), nil
	case jsext.Promise:
		return val.JSValue(), nil
	case []byte:
		var enc = base64.StdEncoding.EncodeToString(val)
		return js.ValueOf(enc), nil
	case int, int64, int32, int16, int8,
		float64, float32,
		uint, uint64, uint32, uint16, uint8, uintptr,
		string, bool:
		// []interface{}, map[string]interface{}: // Removed so we can call jss.ValueOf on a slice or map.

		return js.ValueOf(val), nil
	case func():
		if val == nil {
			return js.Null(), nil
		}
		return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			val()
			return nil
		}).Value, nil
	case func(this js.Value, args []js.Value) interface{}:
		if val == nil {
			return js.Null(), nil
		}
		return js.FuncOf(val).Value, nil
	case Marshaller:
		return val.MarshalJS(), nil
	case ErrorMarshaller:
		var jsValue, err = val.MarshalJS()
		if err != nil {
			return js.Null(), err
		}
		return jsValue, nil
	}
	var valueOf = reflect.ValueOf(f)
	if !valueOf.IsValid() {
		return js.Null(), nil
	}
	var kind = valueOf.Kind()
	return valueOfJS(valueOf, kind)
}

func valueOfJS(valueOf reflect.Value, kind reflect.Kind) (js.Value, error) {
	switch kind {
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return js.ValueOf(valueOf.Int()), nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uintptr:
		return js.ValueOf(valueOf.Uint()), nil
	case reflect.Float64, reflect.Float32:
		return js.ValueOf(valueOf.Float()), nil
	case reflect.String:
		return js.ValueOf(valueOf.String()), nil
	case reflect.Bool:
		return js.ValueOf(valueOf.Bool()), nil
	case reflect.Slice, reflect.Array:
		// Check if bytes
		if valueOf.Type().Elem().Kind() == reflect.Uint8 {
			if !BASE64 {
				var length = valueOf.Len()
				var array = js.Global().Get("Uint8Array").New(length)
				js.CopyBytesToJS(array, valueOf.Bytes())
				return array, nil
			}
			var enc = base64.StdEncoding.EncodeToString(valueOf.Bytes())
			return js.ValueOf(enc), nil
		}
		var length = valueOf.Len()
		var array = js.Global().Get("Array").New(length)
		for i := 0; i < length; i++ {
			var index = valueOf.Index(i)
			if index.Kind() == reflect.Ptr {
				index = index.Elem()
			}

			if !index.IsValid() || !index.CanInterface() {
				return js.Null(), nil
			}

			var v, err = ValueOf(index.Interface())
			if err != nil {
				return js.Null(), err
			}

			array.SetIndex(i, v)
		}
		return js.ValueOf(array), nil
	case reflect.Map:
		var keys = valueOf.MapKeys()
		var object = js.Global().Get("Object").New()
		for _, key := range keys {
			var index = valueOf.MapIndex(key)
			if index.Kind() == reflect.Ptr {
				index = index.Elem()
			}
			if !index.CanInterface() {
				return js.Null(), nil
			}
			var v, err = ValueOf(index.Interface())
			if err != nil {
				return js.Null(), err
			}
			object.Set(key.String(), v)
		}
		return object, nil
	case reflect.Struct:
		if !valueOf.CanInterface() {
			return js.Null(), nil
		}
		if !TINYGO {
			if valueOf.Type().ConvertibleTo(reflect.TypeOf(js.Value{})) {
				var jsValue js.Value
				reflect.ValueOf(&jsValue).Elem().Set(valueOf)
				return jsValue, nil
			}
		}
		var object = js.Global().Get("Object").New()
		var typeOf = valueOf.Type()
		var numField = valueOf.NumField()
		for i := 0; i < numField; i++ {
			var field = typeOf.Field(i)
			var tag, omitEmpty, ok = getStructTag(field, "js", "jsc", "json")
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
				continue
			}

			var v, err = ValueOf(valField.Interface())
			if err != nil {
				return js.Null(), err
			}

			object.Set(tag, v)
		}

		if !TINYGO {
			var numMethod = valueOf.NumMethod()
			for i := 0; i < numMethod; i++ {
				var methodType = valueOf.Type().Method(i)
				var method = valueOf.Method(i)

				if !method.CanInterface() {
					continue
				}

				var v, err = ValueOf(method.Interface())
				if err != nil {
					return js.Null(), err
				}

				object.Set(methodType.Name, v)
			}
		}

		return object, nil
	case reflect.Ptr, reflect.Interface:
		return ValueOf(valueOf.Elem().Interface())
	// Very incompatible with TinyGo...
	case reflect.Func:
		if valueOf.IsNil() {
			return js.Null(), nil
		}
		if TINYGO {
			panic("(reflect.Type).In() not supported in tinygo: cannot convert func to js.Func")
		}
		return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			var in = make([]reflect.Value, len(args))
			for i := range args {
				var arg = args[i]
				var castJS = guessType(arg)
				var jsValueOf = reflect.ValueOf(castJS)
				if jsValueOf.Kind() == reflect.Ptr {
					jsValueOf = jsValueOf.Elem()
				}
				switch {
				case jsValueOf.Type().ConvertibleTo(valueOf.Type().In(i)):
					in[i] = jsValueOf.Convert(valueOf.Type().In(i))
				default:
					in[i] = reflect.ValueOf(castJS)
				}
			}
			var out = valueOf.Call(in)
			if len(out) == 0 {
				return nil
			} else if len(out) == 1 {
				var v, err = ValueOf(out[0].Interface())
				if err != nil {
					return js.Null()
				}
				return v
			}
			var returnValues = make([]interface{}, len(out))
			for i := range out {
				returnValues[i] = out[i].Interface()
			}

			var v, err = ValueOf(returnValues)
			if err != nil {
				return js.Null()
			}

			return v
		}).Value, nil
	default:
		return js.Null(), errs.Error("ValueOf: unsupported type " + kind.String())
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
		for _, part := range parts[1:] {
			if part == "omitempty" {
				omitEmpty = true
				break
			}
		}
	}

	return name, omitEmpty, ok
}

const (
	ErrUndefined errs.Error = "src is null or undefined"
	ErrNil       errs.Error = "dst is nil"
	ErrNotPtr    errs.Error = "dst is not a pointer"
	ErrNotValid  errs.Error = "dst is not a pointer to a struct, map, or slice"
	ErrNotStruct errs.Error = "dst is not a pointer to a struct"
	ErrCannotSet errs.Error = "cannot set dst field"
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

	//	if src.Type() != js.TypeObject {
	//		return ErrNotObject
	//	}

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
		var dstTag, _, dstOk = getStructTag(dstField, "js", "jsc", "json")
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
			if !BASE64 {
				var b = make([]byte, srcVal.Length())
				js.CopyBytesToGo(b, srcVal)
				dstVal.SetBytes(b)
				return nil
			}
			var bytes, err = base64.StdEncoding.DecodeString(srcVal.String())
			if err == nil {
				dstVal.SetBytes(bytes)
			} else {
				var b = make([]byte, srcVal.Length())
				js.CopyBytesToGo(b, srcVal)
			}
			return nil
		}
		var err = scanSlice(srcVal, dstVal.Addr())
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
		i = func(args ...interface{}) js.Value {
			var s = make([]interface{}, len(args))
			for i, arg := range args {
				var v, err = ValueOf(arg)
				if err != nil {
					console.Error(err.Error())
					return js.Null()
				}
				s[i] = v
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

		// check if map value is pointer
		// if it is, set the map value to the pointer
		if dstVal.Type().Elem().Kind() == reflect.Ptr {
			dstVal.SetMapIndex(dstKey.Elem(), dstKeyValue)
			continue
		}

		dstVal.SetMapIndex(dstKey.Elem(), dstKeyValue.Elem())
	}
	return nil
}

func scanSlice(srcVal js.Value, dstVal reflect.Value) error {
	if dstVal.Kind() == reflect.Ptr {
		dstVal = dstVal.Elem()
	}
	var srcLen = srcVal.Length()
	if srcLen == 0 {
		return nil
	}
	if dstVal.IsNil() || dstVal.Len() == 0 {
		// makeslice is implemented in tinygo! :)
		dstVal.Set(reflect.MakeSlice(dstVal.Type(), srcLen, srcLen))
	}
	var elemType = dstVal.Type().Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}
	for i := 0; i < srcLen; i++ {
		var srcElem = srcVal.Index(i)
		var dstElem = reflect.New(elemType)
		var err = scanValue(srcElem, dstElem)
		if err != nil {
			return err
		}

		switch {
		case dstVal.Type().Elem().Kind() != reflect.Ptr:
			dstElem = dstElem.Elem()
		}

		if dstVal.Len() <= i {
			dstVal.Set(reflect.Append(dstVal, dstElem))
			continue
		}

		dstVal.Index(i).Set(dstElem)
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
