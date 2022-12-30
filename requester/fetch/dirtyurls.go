package fetch

import (
	"net/url"
)

type DirtyURLConstraint interface {
	int | int8 | int16 | int32 | int64
	uint | uint8 | uint16 | uint32 | uint64 | bool
	string | float32 | float64 | complex64 | complex128
	[]byte | []rune
}

// Convert a map to a URL encoded string
func ToURLValues[T DirtyURLConstraint](data map[string]T) []byte {
	var d url.Values
	for k, v := range data {
		d.Set(k, toString(v))
	}
	return []byte(d.Encode())
}
