//go:build skipimports
// +build skipimports

package jsc

import (
	"syscall/js"
)

func EncodeBase64[T ~string | ~[]byte](data T) (T, error) {
	var obj = js.ValueOf(data)
	var encoded = js.Global().Get("btoa").Invoke(obj)
	return T(encoded.String()), nil

}

func DecodeBase64[T ~string | ~[]byte](data T) (T, error) {
	var obj = js.Global().Get("atob").Invoke(string(data))
	return T(obj.String()), nil
}
