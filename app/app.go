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

type AppBase struct {
	Root   jsext.Element
	Header components.Component
	Base   *elements.Element
	Footer components.Component
	Loader components.Loader
}

func (a *AppBase) Render() jsext.Element {
	a.Root.Prepend(a.Header.Render())
	a.Root.AppendChild(a.Base.Render())
	a.Root.AppendChild(a.Footer.Render())
	return a.Root
}

func (a *AppBase) Clear() {
	a.Root.InnerHTML("")
}

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
// If id is empty, the application will be initialized on the body.
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
		Base: AppBase{
			Root:   elem,
			Header: nil,
			Footer: nil,
			Loader: loaders.NewLoader(id, loaders.ID_LOADER, true, loaders.LoaderRing),
		},
		Data: make(map[string]interface{}),
	}
	return a
}

// Decide what happens on errors.
func (a *Application) OnError(f func(err error)) {
	var newF = func(err error) {
		f(err)
		a.Base.Render()
	}
	a.onErr = newF
}

// Set the base navbar.
func (a *Application) SetNavbar(navbar components.Component) *Application {
	a.Base.Header = navbar
	return a
}

// Set the base footer.
func (a *Application) SetFooter(footer components.Component) *Application {
	a.Base.Footer = footer
	return a
}

// Set the base loader.
func (a *Application) SetLoader(loader components.Loader) *Application {
	a.Base.Loader = loader
	return a
}

// Set the base style.
func (a *Application) SetStyle(style ...string) *Application {
	a.Base.Base.AttrStyle(style...)
	return a
}

// Set classes on the base element.
func (a *Application) SetClasses(classes ...string) *Application {
	a.Base.Base.AttrClass(classes...)
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

func (a *Application) OnLoad(f func()) *Application {
	a.onLoad = f
	return a
}

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

// Setup application to be ran.
// Return 0 on exit.
func (a *Application) run() int {
	if a.onErr == nil {
		a.OnError(router.DefaultRouterErrorDisplay)
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
	a.Base.Loader.Show()
	go func() {
		f()
		a.Base.Loader.Finalize()
	}()
}

// Register routes to the application.
func (a *Application) Register(name string, path string, callable func(a *Application, v router.Vars, u *url.URL), linkEmpty ...bool) *router.Route {
	var ncall = a.WrapURL(callable)
	if len(linkEmpty) > 0 && linkEmpty[0] {
		ncall = nil
	}
	var route = a.Router.Register(name, path, ncall)
	return route
}

func (a *Application) WrapURL(f func(a *Application, v router.Vars, u *url.URL)) func(v router.Vars, u *url.URL) {
	return func(v router.Vars, u *url.URL) {
		f(a, v, u)
	}
}

// Render a component to the application.
func (a *Application) Render(e ...components.Component) {
	a.render(e...)
}

func (a *Application) Redirect(url string) {
	a.Router.Redirect(url)
}

// Render a component to the application.
func (a *Application) render(e ...components.Component) {
	a.Base.Clear()
	a.Base.Render()
	for _, el := range e {
		a.Base.Base.JSExtElement().AppendChild(el.Render())
	}
}

// InnerHTML sets the inner HTML of the element.
func (a *Application) InnerHTML(html string) *Application {
	a.Base.Clear()
	a.Base.Base.InnerHTML(html)
	a.Base.Render()
	return a
}

// InnerText sets the inner text of the element.
func (a *Application) InnerText(text string) *Application {
	a.Base.Clear()
	a.Base.Base.InnerText(text)
	a.Base.Render()
	return a
}

// SetInnerElement clears the inner HTML and appends the element.
func (a *Application) InnerElement(el jsext.Element) *Application {
	a.Base.Clear()
	a.Base.Base.JSExtElement().AppendChild(el)
	a.Base.Render()
	return a
}

// InnerComponent clears the inner HTML and appends the component.
func (a *Application) InnerComponent(c components.Component) *Application {
	a.Base.Clear()
	a.Base.Base.JSExtElement().AppendChild(c.Render())
	a.Base.Render()
	return a
}

// AppendChild appends a child to the element.
func (a *Application) AppendChild(e *elements.Element) *Application {
	// If footer is not nil, append before it
	//if a.Base.Footer != nil {
	//	var footer, ok = a.Base.Footer.(*elements.Element)
	//	if !ok {
	//		panic("footer is not an element, cannot append before it.")
	//	}
	//	a.Base.InsertBefore(e.Render(), footer.JSExtElement())
	//} else {
	//	a.Base.AppendChild(e.Render())
	//}
	a.Base.Base.JSExtElement().AppendChild(e.Render())
	return a
}

// Exit the application.
func Exit() {
	WAITER <- struct{}{}
	close(WAITER)
}

// Set data on the application.
func (a *Application) Set(k string, v any) {
	a.Data[k] = v
}

// Get data from the application.
func (a *Application) Get(key string) interface{} {
	return a.Data[key]
}

// Get data from the application in the form of an int.
func (a *Application) GetInt(key string) int {
	var data = a.Data[key]
	switch data := data.(type) {
	case int:
		return data
	case int8:
		return int(data)
	case int16:
		return int(data)
	case int32:
		return int(data)
	case int64:
		return int(data)
	}
	return 0
}

// Get data from the application in the form of a uint.
func (a *Application) GetUint(key string) uint {
	var data = a.Data[key]
	switch data := data.(type) {
	case uint:
		return data
	case uint8:
		return uint(data)
	case uint16:
		return uint(data)
	case uint32:
		return uint(data)
	case uint64:
		return uint(data)
	}
	return 0
}

// Get data from the application in the form of a string.
func (a *Application) GetString(key string) string {
	return a.Data[key].(string)
}

// Get data from the application in the form of a bool.
func (a *Application) GetBool(key string) bool {
	return a.Data[key].(bool)
}

// Get data from the application in the form of a float64.
func (a *Application) GetFloat(key string) float64 {
	var data = a.Data[key]
	switch data := data.(type) {
	case float64:
		return data
	case float32:
		return float64(data)
	}
	return 0
}

// Get data from the application in the form of a complex128.
func (a *Application) GetComplex(key string) complex128 {
	var data = a.Data[key]
	switch data := data.(type) {
	case complex128:
		return data
	case complex64:
		return complex128(data)
	}
	return 0
}
