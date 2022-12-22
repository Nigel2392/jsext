//go:build js && wasm && tinygo
// +build js,wasm,tinygo

package app

type Application struct {
	BaseElementID string
	Router        *router.Router
	Navbar        components.Component
	Footer        components.Component
	Loader        components.Loader
	Base          jsext.Element
	onErr         func(err error)
}
