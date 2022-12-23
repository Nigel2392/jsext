//go:build js && wasm
// +build js,wasm

package app

import (
	"net/url"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/components"
	"github.com/Nigel2392/jsext/components/loaders"
	"github.com/Nigel2392/jsext/elements"
	"github.com/Nigel2392/jsext/router"
)

// Preloader to be removed. This should happen automatically from the JS-script.
const JSEXT_PRELOADER_ID = "jsext-preload-container"

// App export to be used for embedding other exports.
var AppExport jsext.Export

// Set application exports.
// Available in javascript console under:
//
//	jsext.App.((defined_methods))
func init() {
	AppExport = jsext.NewExport()
	AppExport.SetFunc("Exit", Exit)
	AppExport.RegisterToExport("App", jsext.JSExt)
}

// Waiter to lock the main thread.
var WAITER = make(chan struct{})

// Initialize a new application.
func App(id string, rt ...*router.Router) *Application {
	// Get the application body
	var elem jsext.Element
	if id == "" {
		elem = jsext.Body
	} else {
		elem = jsext.QuerySelector("#" + id)
	}
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

// Decide what happens on errors.
func (a *Application) OnError(f func(err error)) {
	var newF = func(err error) {
		f(err)
		a.renderBases()
	}
	a.onErr = newF
}

// Set the base navbar.
func (a *Application) SetNavbar(navbar components.Component) *Application {
	a.Navbar = navbar
	return a
}

// Set the base footer.
func (a *Application) SetFooter(footer components.Component) *Application {
	a.Footer = footer
	return a
}

// Set the base loader.
func (a *Application) SetLoader(loader components.Loader) *Application {
	a.Loader = loader
	return a
}

// Set the base style.
func (a *Application) SetStyle(style string) *Application {
	a.Base.SetAttribute("style", style)
	return a
}

// Set classes on the base element.
func (a *Application) SetClasses(class string) *Application {
	a.Base.SetAttribute("class", class)
	return a
}

func (a *Application) SetTitle(title string) *Application {
	jsext.Document.Set("title", title)
	return a
}

// Run the application.
func (a *Application) Run() int {
	return a.run()
}

// Setup application to be ran.
// Return 0 on exit.
func (a *Application) run() int {
	if a.onErr == nil {
		a.OnError(router.DefaultRouterErrorDisplay)
	}
	a.Router.OnError(a.onErr)
	a.Router.Run()
	// Get the preloader, remove it if it exists
	if preloader := jsext.QuerySelector("#" + JSEXT_PRELOADER_ID); !preloader.Value().IsUndefined() {
		preloader.Remove()
	}
	<-WAITER
	return 0
}

// Exit the application.
func (a *Application) Stop() {
	Exit()
}

// Run the application loader for a time consuming function.
func (a *Application) Load(f func()) {
	a.Loader.Run(f)
}

// Register routes to the application.
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

// Render a component to the application.
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

// InnerComponent clears the inner HTML and appends the component.
func (a *Application) InnerComponent(c components.Component) *Application {
	a.Base.InnerHTML("")
	a.Base.AppendChild(c.Render())
	a.renderBases()
	return a
}

// AppendChild appends a child to the element.
func (a *Application) AppendChild(e components.Component) *Application {
	// If footer is not nil, append before it
	if a.Footer != nil {
		var footer, ok = a.Footer.(*elements.Element)
		if !ok {
			panic("footer is not an element, cannot append before it.")
		}
		a.Base.InsertBefore(e.Render(), footer.JSExtElement())
	} else {
		a.Base.AppendChild(e.Render())
	}
	return a
}

// Render application header and footer if defined.
func (a *Application) renderBases() {
	if a.Navbar != nil {
		a.Base.Prepend(a.Navbar.Render())
	}
	if a.Footer != nil {
		a.Base.Append(a.Footer.Render())
	}
}

// Exit the application.
func Exit() {
	WAITER <- struct{}{}
	close(WAITER)
}
