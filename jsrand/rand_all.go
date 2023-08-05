package jsrand

import (
	"encoding/hex"
	"syscall/js"
)

// XorBytes returns a hex encoded string of the XOR of each byte in the string with the previous byte.
func XorBytes(s string) string {
	var b = []byte(s)
	var last = byte(0)
	for i := 0; i < len(b); i++ {
		b[i] ^= last
		last = b[i]
	}
	return hex.EncodeToString(b)
}

func String(n int) string {
	var b = make([]byte, n)
	Reader.Read(b)
	for i := 0; i < n; i++ {
		b[i] = (b[i] % 16) + 97
	}
	return string(b)
}

// Everything below here is taken from crypto/rand_js.go
// This is to avoid a dependency on the crypto package.

// Reader is a global, shared instance of a cryptographically
// secure random number generator.
var Reader interface {
	Read([]byte) (int, error)
}

func init() {
	Reader = &reader{}
}

var jsCrypto = js.Global().Get("crypto")
var uint8Array = js.Global().Get("Uint8Array")

// reader implements a pseudorandom generator
// using JavaScript crypto.getRandomValues method.
// See https://developer.mozilla.org/en-US/docs/Web/API/Crypto/getRandomValues.
type reader struct{}

func (r *reader) Read(b []byte) (int, error) {
	a := uint8Array.New(len(b))
	jsCrypto.Call("getRandomValues", a)
	js.CopyBytesToGo(b, a)
	return len(b), nil
}
