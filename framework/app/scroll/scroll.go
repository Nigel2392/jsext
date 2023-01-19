package scroll

import (
	"strconv"
	"strings"
	"syscall/js"
	"time"
	"unicode"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/framework/components"
	"github.com/Nigel2392/jsext/framework/helpers"
	"github.com/Nigel2392/jsext/framework/requester"
)

// Axis to scroll on.
type Axis int8

// Axis to scroll on.
const (
	ScrollAxisX Axis = 1
	ScrollAxisY Axis = 2
)

// Waiter to lock the main thread.
var waiter = make(chan struct{})

// Page struct.
// Mostly used internally.
type Page struct {
	title  string
	hash   string
	c      components.ComponentWithValue
	onShow func()
	onHide func()
}

// Callback for when the page is being viewed.
func (p *Page) OnShow(cb func()) *Page {
	p.onShow = cb
	return p
}

// Callback for when the page is being hidden.
func (p *Page) OnHide(cb func()) *Page {
	p.onHide = cb
	return p
}

// Application options
type Options struct {
	ScrollAxis    Axis
	ClassPrefix   string
	ScrollThrough bool
}

func (o *Options) setDefaults() {
	if o.ClassPrefix == "" {
		o.ClassPrefix = "jsext-scrollable-app"
	}
	if o.ScrollAxis == 0 {
		o.ScrollAxis = ScrollAxisY
	}
}

// Application struct
type Application struct {
	Loader         components.Loader
	navbar         components.Component
	footer         components.Component
	pages          []*Page
	documentObject jsext.Element
	backgrounds    []*Background
	onPageChange   func(index int)
	currentPage    int
	Options        *Options
	clientFunc     func() *requester.APIClient
	client         *requester.APIClient
}

// Create a new application from options
func App(documentObjectQuerySelector string, options *Options) *Application {
	var object, err = jsext.QuerySelector(documentObjectQuerySelector)
	if err != nil {
		object = jsext.Body
	}
	var s = &Application{
		pages:          make([]*Page, 0),
		documentObject: object,
		Options:        options,
	}
	return s
}

// Set the application loader
func (s *Application) SetLoader(loader components.Loader) {
	s.Loader = loader
}

// Set the application navbar
func (s *Application) SetNavbar(c components.Component) {
	s.navbar = c
}

// Set the application footer
func (s *Application) SetFooter(c components.Component) {
	s.footer = c
}

// Set the application backgrounds
func (s *Application) Backgrounds(t BackgroundType, b ...string) Backgrounds {
	var backgrounds = make([]*Background, len(b))
	for i, background := range b {
		backgrounds[i] = &Background{
			BackgroundType: t,
			Background:     background,
			Gradient: &Gradient{
				Gradients: make([]string, 0),
			},
		}
	}
	s.backgrounds = append(s.backgrounds, backgrounds...)
	return backgrounds
}

// Add a page
func (s *Application) AddPage(title string, c components.ComponentWithValue) *Page {
	var page = &Page{
		title: title,
		hash:  makeSlug(title),
		c:     c,
	}
	s.pages = append(s.pages, page)
	return page
}

// Run the application
func (s *Application) Run() {

	s.Options.setDefaults()

	// Render navbar
	if s.navbar != nil {
		var navbar = s.navbar.Render()
		navbar.ClassList().Add(s.Options.ClassPrefix + "-navbar")
		s.documentObject.AppendChild(navbar)
	}
	// Append the CSS to the document
	var displayDirection string
	var overflowAxis, oppositeOverflowAxis string
	var width string = `calc(100vw * ` + strconv.Itoa(len(s.pages)) + `);`
	var height string = `calc(100vh * ` + strconv.Itoa(len(s.pages)) + `);`
	var axis string
	switch s.Options.ScrollAxis {
	case ScrollAxisX:
		displayDirection = "row"
		overflowAxis = "overflow-x"
		oppositeOverflowAxis = "overflow-y"
		height = "100%"
		axis = "x"
	default:
		displayDirection = "column"
		overflowAxis = "overflow-y"
		oppositeOverflowAxis = "overflow-x"
		width = "100%"
		axis = "y"
	}
	// Styling
	jsext.StyleBlock(s.Options.ClassPrefix+"-navbar-css", func() string {
		var css = `
			* {
				margin: 0;
				padding: 0;
			}
			body {
				height: 100vh;
				width: 100vw;
				overflow: hidden;
			}
			.` + s.Options.ClassPrefix + `-scrollable-page-container {
				width: 100vw;
				height: 100vh;
				` + overflowAxis + `: hidden;
				scroll-behavior: smooth;
				scroll-snap-type: ` + axis + ` mandatory; 
				scroll-snap-stop: always;
				` + oppositeOverflowAxis + `: hidden;
			}
			.` + s.Options.ClassPrefix + `-scrollable-page {
				width: ` + width + `;
				height: ` + height + `;
				display: flex;
				flex-direction: ` + displayDirection + `;
			}
			.` + s.Options.ClassPrefix + `-page {
				scroll-snap-align: center;
				display: flex;
				flex-direction: column;
				align-items: center;
				justify-content: center;
				width: 100vw;
				height: 100vh;
				font-size: 1.5em;
			}
			`
		if s.navbar != nil {
			css += `.` + s.Options.ClassPrefix + `-navbar {
				position: fixed;
				top: 0;
				left: 0;
				right: 0;
				z-index: 1000;
			}`
		}
		if s.footer != nil {
			css += `.` + s.Options.ClassPrefix + `-footer {
				position: fixed;
				bottom: 0;
				left: 0;
				right: 0;
				z-index: 1000;
			}`
		}

		if len(s.backgrounds) > 0 {
			var ct int
			var bg *Background
			var backup = s.backgrounds[0]
			for _, page := range s.pages {
				bg, ct = helpers.GetColor(s.backgrounds, ct, backup)
				css += bg.CSS(`#` + page.hash)
			}
		} else {
			css += (&Background{
				BackgroundType: BackgroundTypeColor,
				Background:     "#333333",
			}).CSS(`.` + s.Options.ClassPrefix + `-page`)
		}

		return css
	}())
	// Create the application elements
	var scrollablePageContainer = jsext.CreateElement("section")
	scrollablePageContainer.ClassList().Add(s.Options.ClassPrefix + "-scrollable-page-container")
	var scrollablePage = jsext.CreateElement("section")
	scrollablePage.ClassList().Add(s.Options.ClassPrefix + "-scrollable-page")
	for _, page := range s.pages {
		var section = jsext.CreateElement("section")
		section.ClassList().Add(s.Options.ClassPrefix + "-page")
		section.Set("id", page.hash)
		var p = page.c.Render()
		p.ClassList().Add(s.Options.ClassPrefix + "-page-content")
		section.AppendChild(p)
		scrollablePage.AppendChild(section)
	}
	scrollablePageContainer.AppendChild(scrollablePage)
	s.documentObject.AppendChild(scrollablePageContainer)

	// Render the footer
	if s.footer != nil {
		var footer = s.footer.Render()
		footer.ClassList().Add(s.Options.ClassPrefix + "-footer")
		s.documentObject.AppendChild(footer)
	}

	// Add the application eventlistener for the arrow keys
	jsext.Document.Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("key").String() == "ArrowRight" {
			s.NextPage()
		} else if args[0].Get("key").String() == "ArrowLeft" {
			s.PreviousPage()
		}
		return nil
	}))

	// Add the application eventlistener for the mouse wheel
	var scrolled = false
	jsext.Document.Call("addEventListener", "wheel", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var dy = args[0].Get("deltaY").Float()
		args[0].Call("preventDefault")
		if !scrolled {
			if dy > 10 {
				s.NextPage()
			} else if dy < -10 {
				s.PreviousPage()
			}
			scrolled = true
			go func() {
				time.Sleep(500 * time.Millisecond)
				scrolled = false
			}()
		}
		return nil
	}), jsext.MapToObject(map[string]any{"passive": false}).Value())

	// Remove the preloader
	const JSEXT_PRELOADER_ID = "jsext-preload-container"
	if preloader, err := jsext.QuerySelector("#" + JSEXT_PRELOADER_ID); err == nil {
		preloader.Remove()
	}

	// Set the initial page
	var hash = jsext.Document.Get("location").Get("hash").String()
	if hash != "" {
		for i, page := range s.pages {
			if page.hash == hash[1:] {
				s.currentPage = i
				break
			}
		}
	}

	s.updatePage()
	<-waiter
}

// Get the page's container.
func (s *Application) containerByIndex(index int) jsext.Element {
	return s.documentObject.Value().QuerySelectorAll("." + s.Options.ClassPrefix + "-page")[index]
}

// Exit the application.
func (s *Application) Close() {
	close(waiter)
}

// Go to the next page
func (s *Application) NextPage() {
	var page = s.pages[s.currentPage]
	if page.onHide != nil {
		page.onHide()
	}
	s.currentPage++
	s.updatePage()
}

// Go to the previous page
func (s *Application) PreviousPage() {
	var page = s.pages[s.currentPage]
	if page.onHide != nil {
		page.onHide()
	}
	s.currentPage--
	s.updatePage()
}

// Update the page
func (s *Application) updatePage() {
	if s.Options.ScrollThrough {
		if s.currentPage >= len(s.pages) {
			s.currentPage = 0
		} else if s.currentPage < 0 {
			s.currentPage = len(s.pages) - 1
		}
	} else {
		if s.currentPage >= len(s.pages) {
			s.currentPage = len(s.pages) - 1
		} else if s.currentPage < 0 {
			s.currentPage = 0
		}
	}
	var page = s.pages[s.currentPage]
	jsext.Document.Set("title", page.title)
	// always push state to /#hash
	js.Global().Get("history").Call("pushState", nil, nil, "#"+page.hash)
	s.containerByIndex(s.currentPage).ScrollIntoView(true)
	if page.onShow != nil {
		page.onShow()
	}
	if s.onPageChange != nil {
		s.onPageChange(s.currentPage)
	}
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

func makeSlug(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	lastLetter := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastLetter = true
		} else {
			if lastLetter {
				b.WriteRune('-')
			}
			lastLetter = false
		}
	}
	return b.String()
}
