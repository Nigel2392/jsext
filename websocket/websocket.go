//go:build js && wasm
// +build js,wasm

package websocket

import (
	"encoding/json"
	"syscall/js"

	"github.com/Nigel2392/jsext/v2"
	"github.com/Nigel2392/jsext/v2/errs"
)

type WebSocket struct {
	value js.Value
	open  bool

	openFuncs    []func(*WebSocket, MessageEvent)
	closeFuncs   []func(*WebSocket, jsext.Event)
	errorFuncs   []func(*WebSocket, jsext.Event)
	messageFuncs []func(*WebSocket, MessageEvent)
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
		return ws, errs.Error("websocket: failed to open")
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
	return errs.Error(reason)
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

// Will convert the bytes to a string, unless connection is binary.
func (w *WebSocket) SendBytes(data []byte) error {
	if !w.open {
		return errs.Error("websocket: not open")
	}

	var arr js.Value

	switch w.BinaryType() {
	case "arraybuffer":
		var uint8Array = js.Global().Get("Uint8Array").New(len(data))
		js.CopyBytesToJS(uint8Array, data)
		arr = uint8Array.Call("buffer")
	default:
		arr = js.ValueOf(string(data))
	}
	w.value.Call("send", arr)

	return nil
}

func (w *WebSocket) SendJSON(v interface{}) error {
	if !w.open {
		return errs.Error("websocket: not open")
	}

	var data, err = json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return w.SendBytes(data)
}

func (w *WebSocket) OnOpen(f func(w *WebSocket, e MessageEvent)) {
	if w.openFuncs == nil {
		w.openFuncs = make([]func(w *WebSocket, e MessageEvent), 0)
	}
	w.openFuncs = append(w.openFuncs, f)
	w.value.Set("onopen", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return nil
		}
		w.open = true
		for _, f := range w.openFuncs {
			f(w, MessageEvent(args[0]))
		}
		return nil
	}))
}

func (w *WebSocket) OnClose(f func(w *WebSocket, e jsext.Event)) {
	if w.closeFuncs == nil {
		w.closeFuncs = make([]func(w *WebSocket, e jsext.Event), 0)
	}
	w.closeFuncs = append(w.closeFuncs, f)
	w.value.Set("onclose", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return nil
		}
		w.open = false
		for _, f := range w.closeFuncs {
			f(w, jsext.Event(args[0]))
		}
		return nil
	}))
}

func (w *WebSocket) OnError(f func(w *WebSocket, e jsext.Event)) {
	if w.errorFuncs == nil {
		w.errorFuncs = make([]func(w *WebSocket, e jsext.Event), 0)
	}
	w.errorFuncs = append(w.errorFuncs, f)
	w.value.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return nil
		}
		w.open = false
		for _, f := range w.errorFuncs {
			f(w, jsext.Event(args[0]))
		}
		return nil
	}))
}

func (w *WebSocket) OnMessage(f func(w *WebSocket, e MessageEvent)) {
	if w.messageFuncs == nil {
		w.messageFuncs = make([]func(w *WebSocket, e MessageEvent), 0)
	}
	w.messageFuncs = append(w.messageFuncs, f)
	w.value.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return nil
		}
		for _, f := range w.messageFuncs {
			f(w, MessageEvent(args[0]))
		}
		return nil
	}))
}
