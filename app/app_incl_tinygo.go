//go:build js && wasm && tinygo
// +build js,wasm,tinygo

package app

import (
	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/components"
	"github.com/Nigel2392/jsext/router"
)

// Main application, holds router and is the core of the
type Application struct {
	BaseElementID string
	Router        *router.Router
	LastURL       *url.URL
	Navbar        components.Component
	Footer        components.Component
	Loader        components.Loader
	Base          jsext.Element
	onErr         func(err error)
	OnLoad        func()
	Data          map[string]interface{}
}
