//go:build js && wasm
// +build js,wasm

package loaders

import (
	"create_elems/jsext"
	"create_elems/jsext/elements"
	"fmt"
	"syscall/js"
)

const (
	ID_LOADER_CONTAINER = "jsext-main-loader-container"
	ID_LOADER           = "jsext-main-loader"
)

type LoaderFunc = func(idContainer, idLoader string) *elements.Element

type Loader struct {
	appendTo       string
	Activated      bool
	jsVal          js.Value
	created        bool
	deleteOnFinish bool
	className      string
	getLoaderElem  LoaderFunc
}

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

func (l *Loader) Stop() {
	l.Delete()
}

func (l *Loader) Run(f func()) {
	l.Show()
	go func() {
		f()
		l.Finalize()
	}()
}

func (l *Loader) Show() {
	if !l.created {
		l.create()
	}
	l.activate()
}

func (l *Loader) Finalize() {
	if l.deleteOnFinish {
		l.Delete()
	} else {
		l.Deactivate()
	}
}

func (l *Loader) create() jsext.Value {
	var loader_container = l.getLoaderElem(ID_LOADER_CONTAINER, ID_LOADER)
	var loader_val = loader_container.RenderTo(l.appendTo).Value()
	l.jsVal = js.Value(loader_val)
	l.created = true
	println("Loader created")
	println(fmt.Sprintf("Loader value: %v", l.jsVal))
	return loader_val
}

func (l *Loader) activate() {
	if !l.Activated {
		l.Activated = true
		l.jsVal.Get("style").Set("display", "block")
	}
}

func (l *Loader) Deactivate() {
	if l.Activated {
		l.Activated = false
		l.jsVal.Get("style").Set("display", "none")
	}
}

func (l *Loader) Delete() {
	l.created = false
	l.jsVal.Call("remove")
}
