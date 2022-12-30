//go:build js && wasm
// +build js,wasm

package components

import (
	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/elements"
)

// Component interface
type Component interface {
	Render() jsext.Element
}

// Component with eventlistener interface
type ComponentWithEventListener interface {
	Component
	AddEventListener(event string, fn func(this jsext.Value, event jsext.Event))
}

// Loader component
type Loader interface {
	Stop()        // Stop the loader.
	Show()        // Show the loader.
	Run(f func()) // Run the function, finalize loader automatically.
	Finalize()    // Finalize loader.
}

type NavBar interface {
	Nav(urls *URLs) Component
}

var AppURLs = &URLs{
	URLs: make([]*URL, 0),
}

type URL struct {
	Hide   chan bool
	Show   chan bool
	Name   string
	URL    string
	Hidden bool
	Elem   *elements.Element
}

type URLs struct {
	URLs []*URL
}

func NewURLs() *URLs {
	return &URLs{
		URLs: make([]*URL, 0),
	}
}

func (u *URLs) Register(name, url string, elem *elements.Element, hidden bool) *URL {
	var URL = &URL{
		Hide:   make(chan bool),
		Show:   make(chan bool),
		Name:   name,
		URL:    url,
		Hidden: hidden,
		Elem:   elem,
	}
	u.URLs = append(u.URLs, URL)
	return URL
}

func (u *URLs) Append(name, url string, hidden bool) {
	u.URLs = append(u.URLs, &URL{
		Hide:   make(chan bool),
		Show:   make(chan bool),
		Name:   name,
		URL:    url,
		Hidden: hidden,
	})
}

func (u *URLs) Get(name string) *URL {
	for _, url := range u.URLs {
		if url.Name == name {
			return url
		}
	}
	return nil
}

func (u *URLs) GetJS(name string) jsext.Element {
	var url = u.Get(name)
	if url != nil {
		return url.Elem.JSExtElement()
	}
	// Return undefined if not found
	return jsext.Element(jsext.Undefined())
}

func (u *URL) HideURL() {
	u.Elem.AttrStyle("display: none")
}

func (u *URL) ShowURL() {
	u.Elem.AttrStyle("display: block")
}

func (url *URL) Run() {
	go func(url *URL) {
		for {
			select {
			case <-url.Hide:
				var jsUrl = jsext.GetElementById("URL-" + url.Name)
				if jsUrl.Value().Truthy() {
					jsUrl.Set("style", "display: none;")
				} else {
					url.Elem.AttrStyle("display: none")
				}
			case <-url.Show:
				var jsUrl = jsext.GetElementById("URL-" + url.Name)
				if jsUrl.Value().Truthy() {
					jsUrl.Set("style", "display: block;")
				} else {
					url.Elem.AttrStyle("display: block")
				}
			}
		}
	}(url)
}
