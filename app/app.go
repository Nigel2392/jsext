//go:build js && wasm
// +build js,wasm

package app

import (
	"net/url"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/components"
	"github.com/Nigel2392/jsext/components/loaders"
	"github.com/Nigel2392/jsext/elements"
	"github.com/Nigel2392/jsext/requester"
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

// Main application, holds router and is the core of the
type Application struct {
	BaseElemSelector string
	Router           *router.Router
	client           *requester.APIClient
	Navbar           components.Component
	Footer           components.Component
	Loader           components.Loader
	Base             jsext.Element
	clientFunc       func() *requester.APIClient
	onErr            func(err error)
	onLoad           func()
	beforeLoad       func()
	Data             DataMap
}

// Initialize a http client with a loader for a new request.
func (a *Application) Client() *requester.APIClient {
	if a.clientFunc != nil {
		a.client = a.clientFunc()
	} else {
		a.client = requester.NewAPIClient()
	}
	a.client.Before(a.Loader.Show)
	a.client.After(func() {
		a.Loader.Finalize()
		a.client = nil
	})
	return a.client
}

// Set the client function.
func (a *Application) SetClientFunc(f func() *requester.APIClient) *Application {
	a.clientFunc = f
	return a
}

// Initialize a new application.
// If id is empty, the application will be initialized on the body.
func App(querySelector string, rt ...*router.Router) *Application {
	// Get the application body
	var elem jsext.Element
	if querySelector == "" {
		elem = jsext.Body
	} else {
		elem = jsext.QuerySelector(querySelector)
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
		BaseElemSelector: querySelector,
		Router:           r,
		Base:             elem,
		Loader:           loaders.NewLoader(querySelector, loaders.ID_LOADER, true, loaders.LoaderRing),
		Data:             make(map[string]interface{}),
	}
	return a
}

// Decide what happens on errors.
func (a *Application) OnError(f func(*Application, error)) {
	var newF = func(err error) {
		f(a, err)
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

// Set title on the document
func (a *Application) SetTitle(title string) *Application {
	jsext.Document.Set("title", title)
	return a
}

// Run the application.
func (a *Application) Run() int {
	return a.run()
}

// Function to be while the application is loading.
func (a *Application) OnLoad(f func()) *Application {
	a.onLoad = f
	return a
}

// Function to be ran before the application is loaded.
func (a *Application) BeforeLoad(f func()) *Application {
	a.beforeLoad = f
	return a
}

// Function to be ran before the router is loaded.
func (a *Application) OnRouterLoad(f func()) *Application {
	a.Router.OnLoad(f)
	return a
}

// Function to be ran before the page is rendered.
func (a *Application) OnPageChange(f func(*Application, router.Vars, *url.URL)) *Application {
	var newF = func(v router.Vars, u *url.URL) {
		f(a, v, u)
	}
	a.Router.OnPageChange(newF)
	return a
}

// Function to be ran after the page is rendered.
func (a *Application) AfterPageChange(f func(*Application, router.Vars, *url.URL)) *Application {
	var newF = func(v router.Vars, u *url.URL) {
		f(a, v, u)
	}
	a.Router.AfterPageChange(newF)
	return a
}

// Setup application to be ran.
// Return 0 on exit.
func (a *Application) run() int {
	if a.onErr == nil {
		a.onErr = func(err error) {
			router.DefaultRouterErrorDisplay(err)
			a.renderBases()
		}
	}
	if a.beforeLoad != nil {
		a.beforeLoad()
	}
	a.Router.OnError(a.onErr)
	a.Router.Run()
	if a.onLoad != nil {
		a.onLoad()
	}
	// Get the preloader, remove it if it exists
	if preloader := jsext.QuerySelector("#" + JSEXT_PRELOADER_ID); preloader.Value().Truthy() {
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
	if a.Loader != nil && f != nil {
		a.Loader.Show()
		go func() {
			f()
			a.Loader.Finalize()
		}()
	}
}

// Register routes to the application.
func (a *Application) Register(name string, path string, callable func(a *Application, v router.Vars, u *url.URL)) *router.Route {
	var ncall func(v router.Vars, u *url.URL)
	if callable != nil {
		ncall = a.WrapURL(callable)
	}
	var route = a.Router.Register(name, path, ncall)
	return route
}

func (a *Application) WrapURL(f func(a *Application, v router.Vars, u *url.URL)) func(v router.Vars, u *url.URL) {
	return func(v router.Vars, u *url.URL) {
		if f != nil {
			f(a, v, u)
		}
	}
}

// Render a component to the application.
func (a *Application) Render(e ...components.Component) {
	a.render(e...)
}

// Redirect to a url.
func (a *Application) Redirect(url string) {
	a.Router.Redirect(url)
}

// Render a component to the application.
func (a *Application) render(e ...components.Component) {
	a.Base.InnerHTML("")
	for _, el := range e {
		a.Base.AppendChild(el.Render())
	}
	a.renderBases()
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
