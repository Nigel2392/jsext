//go:build skipimports
// +build skipimports

package encoding

import (
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/jsc"
)

var (
	jsonGlobal = js.Global().Get("JSON")
	stringify  = jsonGlobal.Get("stringify")
)

func EncodeJSON[T ~string | ~[]byte](data any) (T, error) {
	var json, err = jsc.ValueOf(data)
	if err != nil {
		return T(""), err
	}
	var jsonStr = stringify.Invoke(json)
	return T(jsonStr.String()), nil
}

func DecodeJSON[T ~string | ~[]byte](data T, dst any) error {
	var obj = jsonGlobal.Call("parse", data)
	return jsc.Scan(obj, dst)
}
