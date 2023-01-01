//go:build js && wasm
// +build js,wasm

package jsext

import "syscall/js"

// See js.go.init() for the initialization of this variable onto the global export object.
var Runtime Export = Export(Eval("new EventTarget();"))

// Emit an event on the global Runtime object.
func EventEmit(name string, args ...interface{}) Value {
	var event = js.Global().Get("Event").New(name)
	event.Set("args", args)
	return Runtime.Call("dispatchEvent", event)
}

// Listen for an event on the global Runtime object.
func EventOn(name string, f func(args ...interface{})) Value {
	return Runtime.JSExt().ToElement().AddEventListener(name, func(this Value, event Event) {
		var jsArgs = event.Get("args")
		var args = ArrayToSlice(jsArgs)
		f(args...)
	}).Value()
}

// Listen for multiple events on the global Runtime object.
func EventOnMultiple(f func(args ...interface{}), names ...string) []Value {
	var values []Value = make([]Value, len(names))
	for _, name := range names {
		values = append(values, EventOn(name, f))
	}
	return values
}
