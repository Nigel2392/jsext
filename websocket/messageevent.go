//go:build js && wasm
// +build js,wasm

package websocket

import (
	"syscall/js"
)

type MessageEvent js.Value

func (m MessageEvent) Data() js.Value {
	return js.Value(m).Get("data")
}

func (m MessageEvent) Target() js.Value {
	return js.Value(m).Get("target")
}

func (m MessageEvent) Origin() string {
	return js.Value(m).Get("origin").String()
}

func (m MessageEvent) LastEventId() string {
	return js.Value(m).Get("lastEventId").String()
}

func (m MessageEvent) Source() js.Value {
	return js.Value(m).Get("source")
}

func (m MessageEvent) Ports() js.Value {
	return js.Value(m).Get("ports")
}
