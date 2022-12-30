//go:build js && wasm
// +build js,wasm

package components

import (
	"github.com/Nigel2392/jsext"
)

// Component interface
type Component interface {
	Render() jsext.Element
}

// Component with eventlistener interface
type ComponentWithEventListener interface {
	Component
	AddEventListener(event string, fn func(this jsext.Value, event jsext.Event))
}

// Loader component
type Loader interface {
	Stop()        // Stop the loader.
	Show()        // Show the loader.
	Run(f func()) // Run the function, finalize loader automatically.
	Finalize()    // Finalize loader.
}

// URL for use in components
type URL struct {
	Name string
	Url  string
}
