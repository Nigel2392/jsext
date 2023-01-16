//go:build js && wasm
// +build js,wasm

package components

import (
	"net/url"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/framework/router/routes"
	"github.com/Nigel2392/jsext/framework/router/rterr"
	"github.com/Nigel2392/jsext/framework/router/vars"
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

type Router interface {
	Register(string, string, func(v vars.Vars, u *url.URL)) *routes.Route
	Run()
	OnLoad(f func())
	OnPageChange(func(v vars.Vars, u *url.URL))
	Redirect(string)
	AfterPageChange(f func(vars.Vars, *url.URL))
	OnError(func(err error))
	SkipTrailingSlash()
	NameToTitle(bool)
	Throw(int)
	Use(middleware func(vars.Vars, *url.URL, *routes.Route, rterr.ErrorThrower) bool)
	Error(code int, msg string) rterr.RouterError
}
