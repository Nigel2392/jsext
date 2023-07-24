//go:build js && wasm
// +build js,wasm

package export

import (
	"syscall/js"
)

type simpleError string

func (e simpleError) Error() string {
	return string(e)
}

type Export interface {
	// Value returns the js.Value of the export.
	MarshalJS() js.Value

	// Set a value on the export's js.Value.
	Set(name string, value interface{})

	// Get a value from the export's js.Value.
	Get(name string) js.Value

	// Call a function on the export's js.Value.
	Call(name string, args ...interface{}) js.Value

	// Delete an attribute from the export's js.Value.
	Delete(name string)

	// Remove the export from the DOM window.
	Remove() error
}

// Export resembles an object bound to the DOM window.
// This object contains methods to be called from JS.
type export js.Value

// NewExport returns a new Export.
//
// It will register the export on the DOM window.
func NewExport(exportName string) Export {
	var e = export(js.Global().Get("Object").New())
	e.Set("removeExport", js.FuncOf(func(this js.Value, args []js.Value) any {
		js.Global().Delete(exportName)
		return nil
	}))
	js.Global().Set(exportName, e.MarshalJS())
	return e
}

// NewFromValue returns a new Export from a js.Value.
func NewFromValue(value js.Value) Export {
	return export(value)
}

// Value returns the js.Value of the export.
func (e export) MarshalJS() js.Value {
	return js.Value(e)
}

// Set a value on the export's js.Value.
func (e export) Set(name string, value interface{}) {
	js.Value(e).Set(name, value)
}

// Get a value from the export's js.Value.
func (e export) Get(name string) js.Value {
	return js.Value(e).Get(name)
}

// Call a function on the export's js.Value.
func (e export) Call(name string, args ...interface{}) js.Value {
	return js.Value(e).Call(name, args...)
}

// Delete an attribute from the export's js.Value.
func (e export) Delete(name string) {
	js.Value(e).Delete(name)
}

// Remove the export from the DOM window.
func (e export) Remove() error {
	var removeExport = js.Value(e).Get("removeExport")
	if removeExport.Type() != js.TypeFunction {
		return simpleError("removeExport is not a function")
	}
	removeExport.Invoke()
	return nil
}
