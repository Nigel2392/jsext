//go:build !skipimports
// +build !skipimports

package jsc

import "encoding/base64"

func EncodeBase64[T ~string | ~[]byte](data T) (T, error) {
	return T(base64.StdEncoding.EncodeToString([]byte(data))), nil
}

func DecodeBase64[T ~string | ~[]byte](data T) (T, error) {
	var decoded, err = base64.StdEncoding.DecodeString(string(data))
	return T(decoded), err
}
