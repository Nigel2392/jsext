//go:build !tinygo && js && wasm
// +build !tinygo,js,wasm

package app

import (
	"github.com/Nigel2392/jsext/requester"
	"github.com/Nigel2392/jsext/router"
)

// Main application, holds router and is the core of the
type Application struct {
	BaseElementID string
	Router        *router.Router
	client        *requester.APIClient
	Base          AppBase
	clientFunc    func() *requester.APIClient
	onErr         func(err error)
	onLoad        func()
	Data          map[string]interface{}
}

// Initialize a http client with a loader for a new request.
func (a *Application) Client() *requester.APIClient {
	if a.clientFunc != nil {
		a.client = a.clientFunc()
	} else {
		a.client = requester.NewAPIClient()
	}
	a.client.Before(a.Base.Loader.Show)
	a.client.After(func() {
		a.Base.Loader.Finalize()
		a.client = nil
	})
	return a.client
}

// Set the client function.
func (a *Application) SetClientFunc(f func() *requester.APIClient) *Application {
	a.clientFunc = f
	return a
}
