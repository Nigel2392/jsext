//go:build !js && !wasm
// +build !js,!wasm

package router

import "net/url"

// Router is the main router struct.
type Router struct {
	routes            []*Route
	skipTrailingSlash bool
	onErr             func(err error)
}

// Handle is the main router handler.
// This function is called by the router to match and handle a route.
func (r *Router) Handle(u *url.URL) { //bool
	var rt, vars, ok = r.Match(u.Path)
	if !ok {
		var err = NewError(404, "no route found for path: "+u.Path)
		if r.onErr == nil {
			panic(err)
		} else {
			r.onErr(error(err))
			return //false
		}
	}
	if rt.Callable == nil {
		var err = NewError(500, "no callable for route: "+rt.Name)
		if r.onErr == nil {
			panic(err)
		} else {
			r.onErr(error(err))
			return //false
		}
	}
	rt.Callable(vars, u)
	// return //true
}
