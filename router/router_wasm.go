//go:build js && wasm
// +build js,wasm

package router

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Nigel2392/jsext"
)

// Default router error to be displayed if one occurs.
func DefaultRouterErrorDisplay(err error) {
	var rtErr, ok = err.(RouterError)
	if !ok {
		panic(rtErr)
	}
	var style = jsext.CreateElement("style")
	style.Set("type", "text/css")
	style.Set("id", "jsext-style")
	style.Set("innerHTML", `
	.jsext-overlay { position: absolute; top: 0; left: 0; width: 100%; height: 100%; background-color: rgba(0,0,0,0.5); z-index: 1000; display: flex; justify-content: center; align-items: center; }
	.jsext-modal { background-color: white; padding: 20px; border-radius: 5px; box-shadow: 0 0 10px rgba(0,0,0,0.5); }
	.jsext-modal h1 { margin-top: 5px; font-size: 25px; margin: 0; }
	.jsext-modal p { margin-top: 5px; font-size: 20px; margin: 0; }
	.jsext-modal button { margin-top: 5px; padding: 10px 20px; border: none; border-radius: 5px; background-color: #9200ff; color: white; font-size: 15px; cursor: pointer; }
	.jsext-modal button:hover { background-color: #a200ff; }`)
	var overlay = jsext.CreateElement("div")
	overlay.ClassList().Add("jsext-overlay")
	var modal = jsext.CreateElement("div")
	modal.ClassList().Add("jsext-modal")
	var title = jsext.CreateElement("h1")
	title.ClassList().Add("jsext-modal-title")
	title.InnerHTML("Error")
	var message = jsext.CreateElement("p")
	message.ClassList().Add("jsext-modal-message")
	message.InnerHTML(rtErr.Error())
	var button = jsext.CreateElement("button")
	button.ClassList().Add("jsext-modal-button")
	button.InnerHTML("Close")
	overlay.AddEventListener("click", func(t jsext.Value, event jsext.Event) {
		event.PreventDefault()
		overlay.Remove()
	})
	modal.AppendChild(title)
	modal.AppendChild(message)
	modal.AppendChild(button)
	overlay.AppendChild(modal)
	overlay.AppendChild(style)
	jsext.Body.AppendChild(overlay)
}

// Router is the main router struct.
type Router struct {
	routes            []*Route
	skipTrailingSlash bool
	nameToTitle       bool
	onErr             func(err error)
	onLoad            func()
	onPageChange      func(Vars, *url.URL)
}

// Get a route by index.
func (r *Router) GetIndex(i int) *Route {
	return r.routes[i]
}

// Set on load function.
func (r *Router) OnLoad(f func()) *Router {
	r.onLoad = f
	return r
}

// Set on page change function.
func (r *Router) OnPageChange(f func(Vars, *url.URL)) *Router {
	r.onPageChange = f
	return r
}

// Automatically convert the name of the route to the title of the page.
func (r *Router) NameToTitle(b bool) *Router {
	r.nameToTitle = b
	return r
}

var RT_PREFIX = "router:"
var RT_PREFIX_EXTERNAL = "external:"

// Change the page to the given path.
func (r *Router) changePage(this jsext.Value, event jsext.Event) {
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
	r.HandlePath(path)
}

// Run the router, catch all the events and handle them.
// If a link is clicked, the router will handle it if it has the RT_PREFIX in front.
func (r *Router) Run() {
	jsext.Element(jsext.Document).AddEventListener("click", r.changePage)
	var RouterExport = jsext.NewExport()
	RouterExport.RegisterToExport("Router", jsext.JSExt)
	RouterExport.SetFuncWithArgs("Change", func(this jsext.Value, args jsext.Args) interface{} {
		var path = args[0].String()
		r.HandlePath(path)
		return nil
	})
	RouterExport.SetFuncWithArgs("String", func(this jsext.Value, args jsext.Args) interface{} {
		return r.String()
	})
	RouterExport.SetFuncWithArgs("Routes", func(this jsext.Value, args jsext.Args) interface{} {
		var b strings.Builder
		for _, rt := range r.routes {
			b.WriteString(rt.String())
			b.WriteString(" ")
			b.WriteString(rt.Path)
			b.WriteString("\n")
		}
		return b.String()
	})

	jsext.Element(jsext.Window).AddEventListener("popstate", func(t jsext.Value, event jsext.Event) {
		var path = jsext.Window.Get("location").Get("href").String()
		r.HandlePath(path)
	})
	if r.onLoad != nil {
		r.onLoad()
	}
	var path = jsext.Window.Get("location").Get("href").String()
	r.HandlePath(path)
}

// Handle is the main router handler.
// This function is called by the router to match and handle a route.
func (r *Router) Handle(u *url.URL) {
	go func() {
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
		if r.onPageChange != nil {
			r.onPageChange(vars, u)
		}
		if rt.Callable != nil {
			go rt.Callable(vars, u)
		}
		jsext.Window.Get("history").Call("pushState", nil, "", u.String())
		if r.nameToTitle {
			jsext.Document.Set("title", simpleToTitle(rt.Name))
		}
	}()
}

func (r *Router) Handlef(fmtPath string, args ...any) {
	var path = fmt.Sprintf(fmtPath, args...)
	r.HandlePath(path)
}

func (r *Router) Throw(code int) {
	var err = NewError(code, "error")
	if r.onErr == nil {
		panic(err)
	} else {
		r.onErr(error(err))
	}
}

// Capitalize the first letter of the string.
func simpleToTitle(s string) string {
	var b = []byte(s)
	for i := 0; i < len(b); i++ {
		if b[i] >= 'a' && b[i] <= 'z' {
			b[i] -= 32
			return string(b)
		}
	}
	return string(b)
}
