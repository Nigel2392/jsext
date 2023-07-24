//go:build js && wasm
// +build js,wasm

package jsext

import "syscall/js"

// Wrapper for javascript events to make life easier.
type Event js.Value

// MarshalJS returns the underlying js.Value.
func (e Event) MarshalJS() js.Value {
	return js.Value(e)
}

// Returns the js.Value of the event.
func (e Event) JSValue() js.Value {
	return js.Value(e)
}

// Returns the jsext.Value of the event.
func (e Event) Value() Value {
	return Value(e.JSValue())
}

// Bubbles returns true if the event bubbles.
func (e Event) Bubbles() bool {
	return e.JSValue().Get("bubbles").Bool()
}

// Cancelable returns true if the event is cancelable.
func (e Event) Cancelable() bool {
	return e.JSValue().Get("cancelable").Bool()
}

// Composed returns true if the event is composed.
func (e Event) Composed() bool {
	return e.JSValue().Get("composed").Bool()
}

// CurrentTarget returns the current target of the event.
func (e Event) CurrentTarget() Value {
	return Value(e.JSValue().Get("currentTarget"))
}

// DefaultPrevented returns true if the default action of the event has been prevented.
func (e Event) DefaultPrevented() bool {
	return e.JSValue().Get("defaultPrevented").Bool()
}

// EventPhase returns the event phase of the event.
func (e Event) EventPhase() int {
	return e.JSValue().Get("eventPhase").Int()
}

// IsTrusted returns true if the event is trusted.
func (e Event) IsTrusted() bool {
	return e.JSValue().Get("isTrusted").Bool()
}

// ReturnValue returns the return value of the event.
func (e Event) Target() Value {
	return Value(e.JSValue().Get("target"))
}

// TimeStamp returns the time stamp of the event.
func (e Event) TimeStamp() int {
	return e.JSValue().Get("timeStamp").Int()
}

// Type returns the type of the event.
func (e Event) Type() string {
	return e.JSValue().Get("type").String()
}

// ComposedPath returns the composed path of the event.
func (e Event) ComposedPath() []js.Value {
	composedPath := e.JSValue().Call("composedPath")
	var length = composedPath.Length()
	var path = make([]js.Value, length)
	for i := 0; i < length; i++ {
		path[i] = composedPath.Index(i)
	}
	return path
}

// PreventDefault prevents the default action of the event.
func (e Event) PreventDefault() {
	e.JSValue().Call("preventDefault")
}

// StopImmediatePropagation stops the immediate propagation of the event.
func (e Event) StopImmediatePropagation() {
	e.JSValue().Call("stopImmediatePropagation")
}

// StopPropagation stops the propagation of the event.
func (e Event) StopPropagation() {
	e.JSValue().Call("stopPropagation")
}

// InitEvent initializes the event.
func (e Event) InitEvent(eventType string, bubbles bool, cancelable bool) {
	e.JSValue().Call("initEvent", eventType, bubbles, cancelable)
}

// InitCustomEvent initializes the custom event.
func (e Event) Set(p string, x interface{}) {
	e.JSValue().Set(p, x)
}

// Default js.Value methods
func (e Event) Get(p string) js.Value {
	return e.JSValue().Get(p)
}
func (e Event) Call(m string, args ...interface{}) js.Value {
	args = MarshallableArguments(args...)
	return e.JSValue().Call(m, args...)
}
func (e Event) Delete(p string) {
	e.JSValue().Delete(p)
}
func (e Event) Equal(other js.Value) bool {
	return e.JSValue().Equal(other)
}
