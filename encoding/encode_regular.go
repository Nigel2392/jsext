//go:build !skipimports
// +build !skipimports

package encoding

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
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
