//go:build skipimports
// +build skipimports

package encoding

import (
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/jsc"
)

var (
	jsonGlobal     = js.Global().Get("JSON")
	jsonEncodeFunc = jsonGlobal.Get("stringify")
)

func EncodeJSON[T ~string | ~[]byte](data any) (T, error) {
	return encode[T](data, jsonEncodeFunc)
}

func DecodeJSON[T ~string | ~[]byte](data T, dst any) error {
	var obj = jsonGlobal.Call("parse", string(data))
	return jsc.Scan(obj, dst)
}

func encode[T ~string | ~[]byte](data any, encodeFunc js.Value) (T, error) {
	var obj, err = jsc.ValueOf(data)
	if err != nil {
		return T(""), err
	}
	var encoded = encodeFunc.Invoke(obj)
	return T(encoded.String()), nil
}
