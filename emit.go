//go:build js && wasm
// +build js,wasm

package jsext

import (
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/export"
)

// See js.go.init() for the initialization of this variable onto the global export object.
var Runtime export.Export = export.NewFromValue(Eval("new EventTarget();").Value())

// Emit an event on the global Runtime object.
func EventEmit(name string, args ...interface{}) Value {
	var event = js.Global().Get("Event").New(name)
	event.Set("args", args)
	return Value(Runtime.Call("dispatchEvent", event))
}

// Listen for an event on the global Runtime object.
func EventOn(name string, f func(args ...interface{})) Value {
	return Value(Runtime.Value().Call("addEventListener", name, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var jsArgs = args[0].Get("args")
		var arguments = ArrayToSlice(jsArgs)
		f(arguments...)
		return nil

	})))
}

// Listen for multiple events on the global Runtime object.
func EventOnMultiple(f func(args ...interface{}), names ...string) []Value {
	var values []Value = make([]Value, len(names))
	for _, name := range names {
		values = append(values, EventOn(name, f))
	}
	return values
}
