package jsc

import (
	"encoding/base64"
	"reflect"
	"strings"
	"syscall/js"

	"github.com/Nigel2392/jsext/v2"
	"github.com/Nigel2392/jsext/v2/jse"
)

// Package JSC implements a way to convert javascript objects to go objects, and vice versa.
//
// This package is used to communicate between the frontend and backend.
//
// If the js tag is not specified on struct fields, it will default to using json.
//
// This package uses reflection, thus might cause unforseen errors with TinyGo.

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
		if valueOf.Type().ConvertibleTo(reflect.TypeOf(js.Value{})) {
			var jsValue js.Value
			reflect.ValueOf(&jsValue).Elem().Set(valueOf)
			return jsValue
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
