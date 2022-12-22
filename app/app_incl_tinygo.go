//go:build js && wasm && tinygo
// +build js,wasm,tinygo

package app

import (
	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/components"
	"github.com/Nigel2392/jsext/router"
)

type Application struct {
	BaseElementID string
	Router        *router.Router
	Navbar        components.Component
	Footer        components.Component
	Loader        components.Loader
	Base          jsext.Element
	onErr         func(err error)
}
