//go:build js && wasm
// +build js,wasm

package websocket

import (
	"errors"
	"syscall/js"
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
	ws.value.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ws.open = true
		done <- struct{}{}
		return nil
	}))
	ws.value.Call("addEventListener", "close", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ws.open = false
		return nil
	}))
	ws.value.Call("addEventListener", "error", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ws.open = false
		return nil
	}))

	<-done
	if ws.ReadyState() != SockOpen {
		return nil, errors.New("websocket: failed to open")
	}
	return ws, nil
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
