//go:build js && wasm
// +build js,wasm

package jsext

import "syscall/js"

// Export resembles an object bound to the DOM window.
// This object contains methods to be called from JS.
type Export Value

// NewExport returns a new Export.
func NewExport() Export {
	return Export(js.Global().Get("Object").New())
}

// Value returns the js.Value of the export.
func (e Export) Value() js.Value {
	return js.Value(e)
}

// Set a value on the export's js.Value.
func (e Export) Set(name string, value interface{}) {
	js.Value(e).Set(name, value)
}

// Get a value from the export's js.Value.
func (e Export) Get(name string) Value {
	return Value(js.Value(e).Get(name))
}

// Call a function on the export's js.Value.
func (e Export) Call(name string, args ...interface{}) Value {
	return Value(js.Value(e).Call(name, args...))
}

// SetFunc sets a function on the export's js.Value.
func (e Export) SetFunc(name string, f func()) {
	js.Value(e).Set(name, WrapFunc(f))
}

// SetMultiple sets multiple functions on the export's js.Value.
func (e Export) SetFuncWithArgs(name string, f JSExtFunc) {
	js.Value(e).Set(name, f.ToJSFunc())
}

// SetMultiple sets multiple functions on the export's js.Value.
func (e Export) SetMultipleWithArgs(fns map[string]JSExtFunc) {
	for name, f := range fns {
		e.SetFuncWithArgs(name, f)
	}
}

// Remove a value from the export's js.Value.
func (e Export) Remove(name string) {
	js.Value(e).Delete(name)
}

// Register the export to the global window.
func (e Export) Register(name string) {
	Global.Set(name, e.Value())
}

// RegisterTo registers the export to another jsext.Value.
func (e Export) RegisterTo(name string, to Value) {
	to.Set(name, e.Value())
}

// RegisterToExport registers the export to another export.
func (e Export) RegisterToExport(name string, to Export) {
	to.Set(name, e.Value())
}
