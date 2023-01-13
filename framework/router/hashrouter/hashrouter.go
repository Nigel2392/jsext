//go:build js && wasm
// +build js,wasm

package hashrouter

import (
	"net/url"
	"strings"
	"syscall/js"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/framework/router/routes"
	"github.com/Nigel2392/jsext/framework/router/rterr"
	"github.com/Nigel2392/jsext/framework/router/vars"
)

var RT_PREFIX = "router:"
var RT_PREFIX_EXTERNAL = "external:"

type HashRouter struct {
	routes          map[string]*routes.Route
	nameToTitle     bool
	onErr           func(err error)
	onLoad          func()
	onPageChange    func(vars.Vars, *url.URL)
	afterPageChange func(vars.Vars, *url.URL)
	middlewares     []func(rterr.ErrorThrower) bool
}

// Initialize a new router.
func NewRouter() *HashRouter {
	return &HashRouter{routes: make(map[string]*routes.Route)}
}

// Add a middleware to the router.
func (r *HashRouter) Use(middleware func(rterr.ErrorThrower) bool) *HashRouter {
	r.middlewares = append(r.middlewares, middleware)
	return r
}

// Decide what to do on errors.
func (r *HashRouter) OnError(cb func(err error)) {
	r.onErr = cb
}

// Throw an error in the router with a message.
func (r *HashRouter) Error(code int, msg string) rterr.RouterError {
	var err = rterr.NewError(code, msg)
	if r.onErr == nil {
		panic(err)
	}
	r.onErr(error(err))
	return err
}

// Throw an error in the router with predefined error code messages.
func (r *HashRouter) Throw(code int) {
	var err = rterr.NewError(code)
	if r.onErr == nil {
		panic(err)
	} else {
		r.onErr(error(err))
	}
}

// Set on load function.
func (r *HashRouter) OnLoad(f func()) {
	r.onLoad = f
}

// Set on page change function.
func (r *HashRouter) OnPageChange(f func(vars.Vars, *url.URL)) {
	r.onPageChange = f
}

// Set after page change function.
func (r *HashRouter) AfterPageChange(f func(vars.Vars, *url.URL)) {
	r.afterPageChange = f
}
func (r *HashRouter) NameToTitle(val bool) {
	r.nameToTitle = val
}

// Skiptrailingslash to adhere to the router interface
func (r *HashRouter) SkipTrailingSlash() {
}

// Add a route to the router.
func (r *HashRouter) Register(name, hash string, callable func(v vars.Vars, u *url.URL)) *routes.Route {
	var route = &routes.Route{Name: name, Path: hash, Callable: callable}
	r.routes[hash] = route
	return route
}

// Start the router.
func (r *HashRouter) Run() {
	if r.onLoad != nil {
		r.onLoad()
	}
	r.route()
}

func (r *HashRouter) Match(hash string) (*routes.Route, bool) {
	var rt, ok = r.routes[hash]
	return rt, ok
}

func (r *HashRouter) Handle(hash string) {
	var rt, ok = r.Match(hash)
	if !ok {
		r.Throw(404)
		return
	}
	if r.onPageChange != nil {
		r.onPageChange(nil, nil)
	}
	for _, middleware := range r.middlewares {
		if !middleware(r) {
			return
		}
	}
	if rt.Callable != nil {
		rt.Callable(nil, nil)
	}
	if r.afterPageChange != nil {
		r.afterPageChange(nil, nil)
	}
	if r.nameToTitle {
		jsext.Document.Set("title", rt.Name)
	}

	// Set the location
	js.Global().Get("window").Get("history").Call("pushState", nil, "", hash)
}

func (r *HashRouter) Redirect(hash string) {
	r.Handle(hash)
}

func (r *HashRouter) route() {
	jsext.Element(jsext.Document).AddEventListener("click", r.changePage)
	js.Global().Get("window").Set("onhashchange", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if r.onPageChange != nil {
			r.onPageChange(nil, nil)
		}
		for _, middleware := range r.middlewares {
			if !middleware(r) {
				return nil
			}
		}
		r.Handle(js.Global().Get("window").Get("location").Get("hash").String())
		if r.afterPageChange != nil {
			r.afterPageChange(nil, nil)
		}
		return nil
	}))
	js.Global().Get("window").Call("addEventListener", "popstate", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var event = args[0]
		event.Call("preventDefault")
		hash := js.Global().Get("window").Get("location").Get("hash").String()
		if hash == "" {
			hash = "#"
		}
		return nil
	}))

	hash := js.Global().Get("window").Get("location").Get("hash").String()
	if hash == "" {
		hash = "#"
	}
	r.Handle(hash)
}

func (r *HashRouter) changePage(this jsext.Value, event jsext.Event) {
	// Get the object if it is valid.
	if !event.Value().IsObject() {
		return
	}
	var target = jsext.Element(event.Target())
	if target.Value().IsUndefined() {
		return
	}
	var path = target.Href()
	// Only stop the default action if the link is an internal link
	// Which means it starts with the RT_PREFIX and we need to handle it
	if !strings.HasPrefix(path, RT_PREFIX) {
		if strings.HasPrefix(path, RT_PREFIX_EXTERNAL) {
			path = strings.TrimPrefix(path, RT_PREFIX_EXTERNAL)
			jsext.Window.Get("location").Set("href", path)
		}
		return
	}
	event.PreventDefault()
	path = strings.TrimPrefix(path, RT_PREFIX)
	r.Handle(path)
}
