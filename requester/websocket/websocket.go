//go:build js && wasm
// +build js,wasm

package websocket

import (
	"errors"
	"syscall/js"

	"github.com/Nigel2392/jsext/requester/encoders"
)

type WebSocket struct {
	value js.Value
	open  bool
}

func New(url string, protocols ...string) *WebSocket {
	var arr = js.Global().Get("Array").New(len(protocols))
	for i, protocol := range protocols {
		arr.SetIndex(i, protocol)
	}
	var ws = js.Global().Get("WebSocket").New(url, arr)
	return &WebSocket{value: ws}
}

func (w *WebSocket) Value() js.Value {
	return w.value
}

func Open(url string, protocols ...string) (*WebSocket, error) {
	var ws = New(url, protocols...)
	var done = make(chan struct{})
	var err = make(chan error)

	var openFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ws.open = true
		done <- struct{}{}
		return nil
	})
	ws.value.Call("addEventListener", "open", openFunc)

	var closeFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ws.open = false
		err <- eventToError(args[0])
		return nil
	})
	ws.value.Call("addEventListener", "close", closeFunc)

	select {
	case <-done:
		close(done)
		close(err)
	case e := <-err:
		close(done)
		close(err)
		ws.value.Call("removeEventListener", "close", closeFunc)
		ws.value.Call("removeEventListener", "open", openFunc)
		return ws, e
	}

	ws.value.Call("removeEventListener", "close", closeFunc)
	ws.value.Call("removeEventListener", "open", openFunc)

	if ws.ReadyState() != SockOpen {
		return ws, errors.New("websocket: failed to open")
	}
	return ws, nil
}

func eventToError(event js.Value) error {
	var reason string
	switch event.Get("code").Int() {
	case 1000:
		reason = "Connection finished normally."
	case 1001:
		reason = "Endpoint moving away from the connection."
	case 1002:
		reason = "Protocol error occurred."
	case 1003:
		reason = "Invalid data received."
	case 1004:
		reason = "Reserved."
	case 1005:
		reason = "No status code received."
	case 1006:
		reason = "Connection aborted abnormally."
	case 1007:
		reason = "Invalid data type compared to specified type."
	case 1008:
		reason = "Policy violation."
	case 1009:
		reason = "Message size too large."
	case 1010: // Note that this status code is not used by the server, because it can fail the WebSocket handshake instead.
		reason = "An endpoint (client) is terminating the connection because it has expected the server to negotiate one or more extension, but the server didn't return them in the response message of the WebSocket handshake. <br/> Specifically, the extensions that are needed are: " + event.Get("reason").String()
	case 1011:
		reason = "Unexpected condition prevented the request from being fulfilled."
	case 1015:
		reason = "Certificate error."
	default:
		reason = "Unknown reason"
	}
	return errors.New(reason)
}

func (w *WebSocket) IsOpen() bool {
	return w.open
}

func (w *WebSocket) BinaryType(s ...string) string {
	if len(s) == 0 {
		return w.value.Get("binaryType").String()
	}
	w.value.Set("binaryType", s[0])
	return s[0]
}

func (w *WebSocket) BufferedAmount() int {
	return w.value.Get("bufferedAmount").Int()
}

func (w *WebSocket) Extensions() string {
	return w.value.Get("extensions").String()
}

func (w *WebSocket) Protocol() string {
	return w.value.Get("protocol").String()
}

func (w *WebSocket) ReadyState() ReadyState {
	return ReadyState(w.value.Get("readyState").Int())
}

func (w *WebSocket) URL() string {
	return w.value.Get("url").String()
}

func (w *WebSocket) Close(code ...int) {
	if len(code) == 0 {
		w.value.Call("close")
		return
	}
	w.value.Call("close", code[0])
}

func (w *WebSocket) CloseReasoned(code int, reason string) {
	w.value.Call("close", code, reason)
}

// Allowed inputs:
//   - string
//   - []byte
//   - js.Value
//   - map[string]any
//   - []any
//   - interface{ String() string }
func (w *WebSocket) Send(data interface{}) error {
	if !w.open {
		return errors.New("websocket: not open")
	}
	switch data := data.(type) {
	case string:
		w.value.Call("send", data)
	case []byte:
		var buffer = js.Global().Get("ArrayBuffer").New(len(data))
		var view = js.Global().Get("Uint8Array").New(buffer)
		for i, b := range data {
			view.SetIndex(i, b)
		}
		w.value.Call("send", buffer)
	case js.Value:
		w.value.Call("send", data)
	case map[string]interface{}:
		var d = string(encoders.MarshalMap(data))
		w.value.Call("send", d)
	case []any:
		var d = string(encoders.MarshalList(data))
		w.value.Call("send", d)
	case interface{ String() string }:
		w.value.Call("send", data.String())
	default:
		w.value.Call("send", data)
	}
	return nil
}

func (w *WebSocket) OnOpen(f func()) {
	w.value.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.open = true
		f()
		return nil
	}))
}

func (w *WebSocket) OnClose(f func(js.Value)) {
	w.value.Call("addEventListener", "close", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.open = false
		f(args[0])
		return nil
	}))
}

func (w *WebSocket) OnError(f func(js.Value)) {
	w.value.Call("addEventListener", "error", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.open = false
		f(args[0])
		return nil
	}))
}

func (w *WebSocket) OnMessage(f func(MessageEvent)) {
	w.value.Call("addEventListener", "message", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		f(MessageEvent(args[0]))
		return nil
	}))
}
