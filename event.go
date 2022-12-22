//go:build js && wasm
// +build js,wasm

package jsext

import "syscall/js"

type Event js.Value

func (e Event) JSValue() js.Value {
	return js.Value(e)
}

func (e Event) Value() Value {
	return Value(e.JSValue())
}

func (e Event) Bubbles() bool {
	return e.JSValue().Get("bubbles").Bool()
}

func (e Event) Cancelable() bool {
	return e.JSValue().Get("cancelable").Bool()
}

func (e Event) Composed() bool {
	return e.JSValue().Get("composed").Bool()
}

func (e Event) CurrentTarget() Value {
	return Value(e.JSValue().Get("currentTarget"))
}

func (e Event) DefaultPrevented() bool {
	return e.JSValue().Get("defaultPrevented").Bool()
}

func (e Event) EventPhase() int {
	return e.JSValue().Get("eventPhase").Int()
}

func (e Event) IsTrusted() bool {
	return e.JSValue().Get("isTrusted").Bool()
}

func (e Event) Target() Value {
	return Value(e.JSValue().Get("target"))
}

func (e Event) TimeStamp() int {
	return e.JSValue().Get("timeStamp").Int()
}

func (e Event) Type() string {
	return e.JSValue().Get("type").String()
}
func (e Event) ComposedPath() []js.Value {
	composedPath := e.JSValue().Call("composedPath")
	var length = composedPath.Length()
	var path = make([]js.Value, length)
	for i := 0; i < length; i++ {
		path[i] = composedPath.Index(i)
	}
	return path
}

func (e Event) PreventDefault() {
	e.JSValue().Call("preventDefault")
}

func (e Event) StopImmediatePropagation() {
	e.JSValue().Call("stopImmediatePropagation")
}

func (e Event) StopPropagation() {
	e.JSValue().Call("stopPropagation")
}

func (e Event) InitEvent(eventType string, bubbles bool, cancelable bool) {
	e.JSValue().Call("initEvent", eventType, bubbles, cancelable)
}
func (e Event) Set(p string, x interface{}) {
	e.JSValue().Set(p, x)
}

func (e Event) Get(p string) js.Value {
	return e.JSValue().Get(p)
}
func (e Event) Call(m string, args ...interface{}) js.Value {
	return e.JSValue().Call(m, args...)
}

func (e Event) Delete(p string) {
	e.JSValue().Delete(p)
}

func (e Event) Equal(other js.Value) bool {
	return e.JSValue().Equal(other)
}
