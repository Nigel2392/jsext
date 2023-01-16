package router

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/Nigel2392/jsext/framework/router/routes"
	"github.com/Nigel2392/jsext/framework/router/rterr"
	"github.com/Nigel2392/jsext/framework/router/vars"
)

// Router is the main router struct.
type Router struct {
	routes            []*routes.Route
	skipTrailingSlash bool
	nameToTitle       bool
	onErr             func(err error)
	onLoad            func()
	onPageChange      func(vars.Vars, *url.URL)
	afterPageChange   func(vars.Vars, *url.URL)
	middlewares       []func(vars.Vars, *url.URL, *routes.Route, rterr.ErrorThrower) bool
}

// Initialize a new router.
func NewRouter() *Router {
	return &Router{routes: make([]*routes.Route, 0)}
}

// SkipTrailingSlash will skip the trailing slash in the path.
func (r *Router) SkipTrailingSlash() {
	r.skipTrailingSlash = true
}

// Add a middleware to the router.
func (r *Router) Use(middleware func(vars.Vars, *url.URL, *routes.Route, rterr.ErrorThrower) bool) {
	r.middlewares = append(r.middlewares, middleware)
}

// Decide what to do on errors.
func (r *Router) OnError(cb func(err error)) {
	r.onErr = cb
}

// Throw an error in the router with a message.
func (r *Router) Error(code int, msg string) rterr.RouterError {
	var err = rterr.NewError(code, msg)
	if r.onErr == nil {
		panic(err)
	}
	r.onErr(error(err))
	return err
}

// Throw an error in the router with predefined error code messages.
func (r *Router) Throw(code int) {
	var err = rterr.NewError(code)
	if r.onErr == nil {
		panic(err)
	} else {
		r.onErr(error(err))
	}
}

// Display nicely formatted URLs
func (r *Router) String() string {
	var sb = &strings.Builder{}
	sb.WriteString("Errorfunc defined: " + strconv.FormatBool(r.onErr != nil) + "\n")
	sb.WriteString("Skip trailing slash: " + strconv.FormatBool(r.skipTrailingSlash) + "\n")
	var level = 0
	for _, route := range r.routes {
		route.StringIndent(sb, level)
	}
	return sb.String()
}

// Match a raw path.
func (r *Router) Match(path string) (*routes.Route, vars.Vars, bool) {
	if r.skipTrailingSlash && len(path) > 1 {
		path = strings.TrimSuffix(path, "/")
	}
	for _, route := range r.routes {
		if match, rt, vars := route.Match(path); match {
			return rt, vars, match
		}
	}
	return &routes.Route{
		Internal_name: "404",
		Path:          path,
	}, nil, false
}

// Get a route by name.
// If it does not directly exist, it will search in the subroutes.
func (r *Router) GetRoute(name string) *routes.Route {
	for _, route := range r.routes {
		if route.Internal_name == name {
			return route
		}
	}
	for _, route := range r.routes {
		var rt = route.GetRoute(name)
		if rt != nil {
			return rt
		}
	}
	return nil
}

// Register a new route.
// If the route name already exists, it will panic.
func (r *Router) Register(name, path string, callable func(v vars.Vars, u *url.URL)) *routes.Route {
	if r.skipTrailingSlash && len(path) > 1 {
		path = strings.TrimSuffix(path, "/")
	}
	for _, route := range r.routes {
		if route.Internal_name == name {
			panic("Router [500] route name already exists: " + name)
		}
	}
	var route = &routes.Route{Name: name, Internal_name: name, Path: path, Callable: callable, SkipTrailingSlash: r.skipTrailingSlash}
	r.routes = append(r.routes, route)
	return route
}

// Handle a path.
func (r *Router) HandlePath(path string) {
	var u, err = url.Parse(path)
	if err != nil {
		r.onErr(err)
		return
	}
	r.Handle(u)
}

// Handle a path in the form of a redirect. NYI.
func (r *Router) Redirect(path string) {
	r.HandlePath(path)
}

// Redirect to a route by name.
func (r *Router) RedirectNamed(name string, varMap vars.Vars) {
	var route = r.GetRoute(name)
	if route == nil {
		r.Error(404, "Route not found: "+name)
		return
	}
	if varMap == nil {
		varMap = make(vars.Vars)
	}
	if route.Callable != nil {
		go route.Callable(varMap, nil)
	}
}
