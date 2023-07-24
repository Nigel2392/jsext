package jsc

import (
	"fmt"
	"reflect"
	"syscall/js"
	"unsafe"

	"github.com/Nigel2392/jsext/v2"
	"github.com/Nigel2392/jsext/v2/errs"
)

// Convert a js.Value to a map[string]T.
type ObjectConstraints interface {
	string | int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float64 | float32 | bool | js.Value |
		[]byte | []interface{} | map[string]interface{} | interface{}
}

// Convert a js.Value array to a []interface{}.
func JSToSlice[T ObjectConstraints](arr js.Value) (array []T, errorCasting error) {
	return array, nil
}

func IsArray(value js.Value) bool {
	return IsSlice(value) || IsTypedArray(value)
}

func IsSlice(value js.Value) bool {
	return js.Global().Get("Array").Call("isArray", value).Bool()
}

func IsTypedArray(value js.Value) bool {
	return js.Global().Get("ArrayBuffer").Call("isView", value).Bool() && !IsDataView(value)
}

func IsDataView(value js.Value) bool {
	return value.InstanceOf(js.Global().Get("DataView"))
}

// Cast a value to the correct type, return an error
// if the javascript value is not the correct type.
//
// Supported types are:
//
//   - string 							   (js.TypeString)
//   - int, int8, int16, int32, int64 	   (js.TypeNumber)
//   - uint, uint8, uint16, uint32, uint64 (js.TypeNumber)
//   - float32, float64 				   (js.TypeNumber)
//   - bool 							   (js.TypeBoolean)
//   - js.Value 						   (js.TypeObject)
//   - []byte 							   (js.TypeObject)
//   - []interface{} 					   (js.TypeObject)
//   - map[string]interface{} 			   (js.TypeObject)
//   - interface{} 						   (we don't know the type, but we can try to guess it.)
func Cast[T ObjectConstraints](value js.Value) (T, error) {
	var preTypeOf T
	var typeOf = interface{}(preTypeOf)

	if value.Type() == js.TypeNull {
		return preTypeOf, errs.Error("cannot cast null")
	} else if value.Type() == js.TypeUndefined {
		return preTypeOf, errs.Error("cannot cast undefined")
	}

	// Do some type checking on the javascript values.
	//
	// This is to try and prevent panics when casting.
	//
	// this is because there is no quick way to check if a js.Value is an array.
	switch typeOf.(type) {
	case js.Value:
		return *(*T)(unsafe.Pointer(&value)), nil
	case string:
		if value.Type() != js.TypeString {
			return preTypeOf, errs.Error("cannot cast to string")
		}
	case int, int8, int16, int32, int64, float32, float64, uint, uint8, uint16, uint32, uint64:
		if value.Type() != js.TypeNumber {
			return preTypeOf, errs.Error("cannot cast to number")
		}
	case bool:
		if value.Type() != js.TypeBoolean {
			return preTypeOf, errs.Error("cannot cast to boolean")
		}
	case map[string]interface{}:
		if value.Type() != js.TypeObject {
			return preTypeOf, errs.Error("cannot cast to object")
		}
	case []byte:
		if !IsTypedArray(value) {
			return preTypeOf, errs.Error("cannot cast to byte array")
		}
	case []interface{}:
		if !IsSlice(value) {
			return preTypeOf, errs.Error("cannot cast to slice")
		}
	}

	// Cast the javascript values to the correct type.
	//
	// This is done by using unsafe.Pointer to convert the value to the correct type
	//
	// and avoiding the extra interface{} conversion.
	switch typeOf.(type) {
	case jsext.Unmarshaller:
		var err = typeOf.(jsext.Unmarshaller).UnmarshalJS(value)
		if err != nil {
			return preTypeOf, err
		}
		return typeOf.(T), nil
	case string:
		var s = value.String()
		return *(*T)(unsafe.Pointer(&s)), nil
	case int:
		var i = value.Int()
		return *(*T)(unsafe.Pointer(&i)), nil
	case int8:
		var i = int8(value.Int())
		return *(*T)(unsafe.Pointer(&i)), nil
	case int16:
		var i = int16(value.Int())
		return *(*T)(unsafe.Pointer(&i)), nil
	case int32:
		var i = int32(value.Int())
		return *(*T)(unsafe.Pointer(&i)), nil
	case int64:
		var i = value.Int()
		return *(*T)(unsafe.Pointer(&i)), nil
	case float32:
		var f = float32(value.Float())
		return *(*T)(unsafe.Pointer(&f)), nil
	case float64:
		var f = value.Float()
		return *(*T)(unsafe.Pointer(&f)), nil
	case bool:
		var b = value.Bool()
		return *(*T)(unsafe.Pointer(&b)), nil
	case []byte:
		var b = make([]byte, value.Length())
		js.CopyBytesToGo(b, value)
		return *(*T)(unsafe.Pointer(&b)), nil
	case []interface{}:
		var a, err = castArray[interface{}](value)
		if err != nil {
			return preTypeOf, err
		}
		return *(*T)(unsafe.Pointer(&a)), nil
	case map[string]interface{}:
		var err error
		var m = make(map[string]interface{})
		var keys = js.Global().Get("Object").Call("keys", value)
		for i := 0; i < keys.Length(); i++ {
			var key = keys.Index(i).String()
			var value = value.Get(key)
			m[key], err = Cast[any](value)
			if err != nil {
				return preTypeOf, err
			}
		}
		return *(*T)(unsafe.Pointer(&m)), nil
	case uint:
		var u = uint(value.Int())
		return *(*T)(unsafe.Pointer(&u)), nil
	case uint8:
		var u = uint8(value.Int())
		return *(*T)(unsafe.Pointer(&u)), nil
	case uint16:
		var u = uint16(value.Int())
		return *(*T)(unsafe.Pointer(&u)), nil
	case uint32:
		var u = uint32(value.Int())
		return *(*T)(unsafe.Pointer(&u)), nil
	case uint64:
		var u = uint64(value.Int())
		return *(*T)(unsafe.Pointer(&u)), nil
	case js.Value:
		return *(*T)(unsafe.Pointer(&value)), nil
	}
	// If the value is an interface, we cannot do switch case type checking.
	if valueIsNullInterface[T]() {
		var v, err = guessConvert(value)
		if err != nil {
			return preTypeOf, err
		}
		return *(*T)(unsafe.Pointer(&v)), nil
	}
	return preTypeOf, fmt.Errorf("cannot cast to %T", typeOf)
}

func castArray[T any](value js.Value) ([]T, error) {
	if !IsArray(value) {
		return []T{}, errs.Error("cannot cast to array")
	}
	var a = make([]T, value.Length())
	var (
		v   T
		err error
	)
	for i := 0; i < value.Length(); i++ {
		v, err = Cast[T](value.Index(i))
		if err != nil {
			return []T{}, err
		}
		a[i] = v
	}
	return a, nil
}

// guessConvert tries to guess the type of the javascript value and convert it to the correct type.
func guessConvert(value js.Value) (interface{}, error) {
	switch value.Type() {
	case js.TypeNull, js.TypeUndefined:
		return nil, nil
	case js.TypeBoolean:
		return value.Bool(), nil
	case js.TypeNumber:
		//// Check if the number is an integer.
		//if value.Float() == float64(value.Int()) {
		//	return value.Int(), nil
		//}
		return value.Float(), nil
	case js.TypeString:
		return value.String(), nil
	case js.TypeSymbol:
		return value, nil
	case js.TypeObject:
		if IsArray(value) {
			return castArray[interface{}](value)
		} else if IsDataView(value) {
			return value, nil
		} else if IsTypedArray(value) {
			var b = make([]byte, value.Length())
			js.CopyBytesToGo(b, value)
			return b, nil
		} else {
			return Cast[map[string]interface{}](value)
		}
	case js.TypeFunction:
		return value, nil
	}
	return nil, errs.Error("cannot convert to type")
}

func valueIsNullInterface[T any]() bool {
	var v T
	return reflect.TypeOf(v) == reflect.TypeOf(*new(any))
}
