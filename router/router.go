package router

import (
	"net/url"
	"strconv"
	"strings"
)

func NewRouter() *Router {
	return &Router{routes: make([]*Route, 0)}
}

func (r *Router) SkipTrailingSlash() *Router {
	r.skipTrailingSlash = true
	return r
}

func (r *Router) OnError(cb func(err error)) *Router {
	r.onErr = cb
	return r
}

func (r *Router) Error(code int, msg string) RouterError {
	var err = NewError(code, msg)
	if r.onErr == nil {
		panic(err)
	}
	r.onErr(error(err))
	return err
}

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
		Name: "404",
		Path: path,
	}, nil, false
}

func (r *Router) GetRoute(name string) *Route {
	for _, route := range r.routes {
		if route.Name == name {
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

func (r *Router) Register(name, path string, callable func(v Vars, u *url.URL)) *Route {
	if r.skipTrailingSlash && len(path) > 1 {
		path = strings.TrimSuffix(path, "/")
	}
	for _, route := range r.routes {
		if route.Name == name {
			panic("Router [500] route name already exists: " + name)
		}
	}
	var route = &Route{Name: name, Path: path, Callable: callable, skipTrailingSlash: r.skipTrailingSlash}
	r.routes = append(r.routes, route)
	return route
}

func (r *Router) HandlePath(path string) {
	var u, err = url.Parse(path)
	if err != nil {
		r.onErr(err)
		return
	}
	r.Handle(u)
}

func (r *Router) Redirect(path string) {
	r.HandlePath(path)
}

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
		route.Callable(vars, nil)
	}
}
