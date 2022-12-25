package router

import (
	"net/url"
	"strconv"
	"strings"
)

// Initialize a new router.
func NewRouter() *Router {
	return &Router{routes: make([]*Route, 0)}
}

// SkipTrailingSlash will skip the trailing slash in the path.
func (r *Router) SkipTrailingSlash() *Router {
	r.skipTrailingSlash = true
	return r
}

// Decide what to do on errors.
func (r *Router) OnError(cb func(err error)) *Router {
	r.onErr = cb
	return r
}

// Throw an error in the router.
func (r *Router) Error(code int, msg string) RouterError {
	var err = NewError(code, msg)
	if r.onErr == nil {
		panic(err)
	}
	r.onErr(error(err))
	return err
}

// Display nicely formatted URLs
func (r *Router) String() string {
	var sb strings.Builder
	for _, route := range r.routes {
		sb.WriteString(route.String())
		sb.WriteString("\n")
	}
	sb.WriteString("Errorfunc defined: " + strconv.FormatBool(r.onErr != nil) + "\n")
	sb.WriteString("Skip trailing slash: " + strconv.FormatBool(r.skipTrailingSlash) + "\n")

	return sb.String()
}

// Match a raw path.
func (r *Router) Match(path string) (*Route, Vars, bool) {
	if r.skipTrailingSlash && len(path) > 1 {
		path = strings.TrimSuffix(path, "/")
	}
	for _, route := range r.routes {
		if match, rt, vars := route.Match(path); match {
			return rt, vars, match
		}
	}
	return &Route{
		name: "404",
		Path: path,
	}, nil, false
}

// Get a route by name.
// If it does not directly exist, it will search in the subroutes.
func (r *Router) GetRoute(name string) *Route {
	for _, route := range r.routes {
		if route.name == name {
			return route
		}
	}
	for _, route := range r.routes {
		var rt = route.getRoute(name)
		if rt != nil {
			return rt
		}
	}
	return nil
}

// Register a new route.
// If the route name already exists, it will panic.
func (r *Router) Register(name, path string, callable func(v Vars, u *url.URL)) *Route {
	if r.skipTrailingSlash && len(path) > 1 {
		path = strings.TrimSuffix(path, "/")
	}
	for _, route := range r.routes {
		if route.name == name {
			panic("Router [500] route name already exists: " + name)
		}
	}
	var route = &Route{Name: name, name: name, Path: path, Callable: callable, skipTrailingSlash: r.skipTrailingSlash}
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
func (r *Router) RedirectNamed(name string, vars Vars) {
	var route = r.GetRoute(name)
	if route == nil {
		r.Error(404, "Route not found: "+name)
		return
	}
	if vars == nil {
		vars = make(Vars)
	}
	if route.Callable != nil {
		go route.Callable(vars, nil)
	}
}
