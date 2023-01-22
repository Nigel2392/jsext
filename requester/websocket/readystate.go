//go:build js && wasm
// +build js,wasm

package websocket

type ReadyState int

func (r ReadyState) String() string {
	switch r {
	case SockConnecting:
		return "Connecting"
	case SockOpen:
		return "Open"
	case SockClosing:
		return "Closing"
	case SockClosed:
		return "Closed"
	}
	return "Unknown"
}

func (r ReadyState) Is(c int) bool {
	return r == ReadyState(c)
}

const (
	SockConnecting ReadyState = 0
	SockOpen       ReadyState = 1
	SockClosing    ReadyState = 2
	SockClosed     ReadyState = 3
)
