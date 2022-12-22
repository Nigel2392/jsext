//go:build js && wasm
// +build js,wasm

package app

import (
	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/components"
	"github.com/Nigel2392/jsext/components/loaders"
	"github.com/Nigel2392/jsext/router"
	"net/url"
)

var AppExport jsext.Export

func init() {
	AppExport = jsext.NewExport()
	AppExport.SetFunc("Exit", Exit)
	AppExport.RegisterToExport("App", jsext.JSExt)
}

var WAITER = make(chan struct{})

func App(id string, rt ...*router.Router) *Application {
	// Get the application body
	if id == "" {
		id = "app"
	}
	var elem = jsext.QuerySelector("#" + id)
	// Get the application router
	var r *router.Router
	if len(rt) > 0 {
		r = rt[0]
	} else {
		r = router.NewRouter()
		r.SkipTrailingSlash()
		r.NameToTitle(true)
	}
	// Return the application
	var a = &Application{
		BaseElementID: id,
		Router:        r,
		Base:          elem,
		Loader:        loaders.NewLoader(id, loaders.ID_LOADER, true, loaders.LoaderRing),
	}
	return a
}

func (a *Application) OnError(f func(err error)) {
	var newF = func(err error) {
		f(err)
		a.renderBases()
	}
	a.onErr = newF
}

func (a *Application) SetNavbar(navbar components.Component) *Application {
	a.Navbar = navbar
	return a
}

func (a *Application) SetFooter(footer components.Component) *Application {
	a.Footer = footer
	return a
}

func (a *Application) SetLoader(loader components.Loader) *Application {
	a.Loader = loader
	return a
}

func (a *Application) SetStyle(style string) *Application {
	a.Base.SetAttribute("style", style)
	return a
}

func (a *Application) SetClass(class string) *Application {
	a.Base.SetAttribute("class", class)
	return a
}

func (a *Application) Run() int {
	return a.run()
}

func (a *Application) run() int {
	if a.onErr == nil {
		a.OnError(router.DefaultRouterErrorDisplay)
	}
	a.Router.OnError(a.onErr)
	a.Router.Run()
	<-WAITER
	return 0
}

func (a *Application) Stop() {
	Exit()
}

func (a *Application) Load(f func()) {
	a.Loader.Run(f)
}

func (a *Application) Register(name string, path string, callable func(a *Application, v router.Vars, u *url.URL), linkEmpty ...bool) *Application {
	var ncall = func(v router.Vars, u *url.URL) {
		callable(a, v, u)
	}
	if len(linkEmpty) > 0 && linkEmpty[0] {
		ncall = nil
	}
	a.Router.Register(name, path, ncall)
	return a
}

func (a *Application) Render(e components.Component) {
	a.InnerElement(e.Render())
}

// InnerHTML sets the inner HTML of the element.
func (a *Application) InnerHTML(html string) *Application {
	a.Base.InnerHTML(html)
	a.renderBases()
	return a
}

// InnerText sets the inner text of the element.
func (a *Application) InnerText(text string) *Application {
	a.Base.InnerText(text)
	a.renderBases()
	return a
}

// SetInnerElement clears the inner HTML and appends the element.
func (a *Application) InnerElement(el jsext.Element) *Application {
	a.Base.InnerHTML("")
	a.Base.AppendChild(el)
	a.renderBases()
	return a
}

func (a *Application) renderBases() {
	if a.Navbar != nil {
		a.Base.Prepend(a.Navbar.Render())
	}
	if a.Footer != nil {
		a.Base.Append(a.Footer.Render())
	}
}

func Exit() {
	WAITER <- struct{}{}
	close(WAITER)
}
