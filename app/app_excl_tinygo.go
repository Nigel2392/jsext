//go:build !tinygo
// +build !tinygo

package app

import (
	"create_elems/jsext"
	"create_elems/jsext/components"
	"create_elems/jsext/requester"
	"create_elems/jsext/router"
)

type Application struct {
	BaseElementID string
	Router        *router.Router
	client        *requester.APIClient
	Navbar        components.Component
	Footer        components.Component
	Loader        components.Loader
	Base          jsext.Element
	onErr         func(err error)
}

func (a *Application) Client() *requester.APIClient {
	a.client = requester.NewAPIClient()
	a.client.Before(a.Loader.Show)
	a.client.After(func() {
		a.Loader.Finalize()
		a.client = nil
	})
	return a.client
}
