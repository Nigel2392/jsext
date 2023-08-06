//go:build !skipimports
// +build !skipimports

package encoding

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
)

func EncodeJSON[T ~string | ~[]byte](data any) (T, error) {
	var b = new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(data); err != nil {
		return T(""), err
	}
	return T(b.Bytes()), nil
}

func DecodeJSON[T ~string | ~[]byte](data T, dst any) error {
	return json.NewDecoder(
		bytes.NewReader([]byte(data)),
	).Decode(dst)
}

func EncodeBase64[T ~string | ~[]byte](data T) (T, error) {
	return T(base64.StdEncoding.EncodeToString([]byte(data))), nil
}

func DecodeBase64[T ~string | ~[]byte](data T) (T, error) {
	var decoded, err = base64.StdEncoding.DecodeString(string(data))
	return T(decoded), err
}
