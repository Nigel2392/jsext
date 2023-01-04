//go:build js && wasm
// +build js,wasm

package loaders

import (
	"syscall/js"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/elements"
)

const (
	ID_LOADER_CONTAINER = "jsext-main-loader-container"
	ID_LOADER           = "jsext-main-loader"
)

// LoaderFunc is a function that returns a loader element.
type LoaderFunc = func(idContainer, idLoader string) *elements.Element

// Loader struct to be used for all loader elements.
type Loader struct {
	appendTo       string
	Activated      bool
	jsVal          js.Value
	created        bool
	deleteOnFinish bool
	className      string
	getLoaderElem  LoaderFunc
}

// Returns a new loader.
func NewLoader(appendTo string, className string, deleteOnFinish bool, f ...LoaderFunc) *Loader {
	var loadFunc = LoaderRing
	if len(f) > 0 {
		loadFunc = f[0]
	}
	return &Loader{
		appendTo:       appendTo,
		className:      className,
		getLoaderElem:  loadFunc,
		deleteOnFinish: deleteOnFinish,
	}
}

// Stop the loader.
func (l *Loader) Stop() {
	l.Finalize()
}

// Run the loader.
func (l *Loader) Run(f func()) {
	l.Show()
	go func() {
		f()
		l.Finalize()
	}()
}

// Show the loader
func (l *Loader) Show() {
	if !l.created {
		l.create()
	}
	l.activate()
}

// Delete or deactivate depending on the deleteOnFinish flag.
func (l *Loader) Finalize() {
	if l.deleteOnFinish {
		l.Delete()
	} else {
		l.Deactivate()
	}
}

// Create the loader element.
func (l *Loader) create() jsext.Value {
	var loader_container = l.getLoaderElem(ID_LOADER_CONTAINER, ID_LOADER)
	var loader_val = loader_container.RenderTo(l.appendTo).Value()
	l.jsVal = js.Value(loader_val)
	l.created = true
	return loader_val
}

// Activate the loader.
func (l *Loader) activate() {
	if !l.Activated {
		l.Activated = true
		l.jsVal.Get("style").Set("display", "block")
	}
}

// Deactivate the loader.
func (l *Loader) Deactivate() {
	if l.Activated {
		l.Activated = false
		l.jsVal.Get("style").Set("display", "none")
	}
}

// Delete the loader.
func (l *Loader) Delete() {
	l.created = false
	l.jsVal.Call("remove")
}
