package encoders

import (
	"net/url"
)

// Convert a map to a URL encoded string
// Can only convert primitive types.
func ToURLValues(data map[string]any) []byte {
	var d url.Values
	for k, v := range data {
		d.Set(k, toString(v))
	}
	return []byte(d.Encode())
}
