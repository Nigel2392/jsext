package elements

import (
	"github.com/Nigel2392/jsext/router"
	"strings"
	"sync"
)

// A url map element
type URLs struct {
	order []string
	urls  map[string]*Element
	mu    *sync.Mutex
}

// Create a new url map
func NewURLs() *URLs {
	return &URLs{
		urls:  make(map[string]*Element),
		order: make([]string, 0),
		mu:    &sync.Mutex{},
	}
}

// Get a url element from the map
func (u *URLs) Get(key string) *Element {
	key = strings.ToUpper(key)
	return u.urls[key]
}

// Set a url element in the map
func (u *URLs) Set(key string, value *Element, external ...bool) {
	if len(external) == 0 || len(external) > 0 && !external[0] {
		var href = value.GetAttr("href")
		// if len(href) > 6 && href[:6] != router.RT_PREFIX {
		value.Delete("href")
		value.AttrHref(router.RT_PREFIX + href)
		// }
	}
	if value.Text == "" {
		value.Text = key
	}
	key = strings.ToUpper(key)
	if _, ok := u.urls[key]; ok {
		panic("URL already exists: " + key)
	}
	u.mu.Lock()
	u.urls[key] = value
	u.order = append(u.order, key)
	u.mu.Unlock()
}

// Delete a url element from the map
func (u *URLs) Delete(key string) {
	key = strings.ToUpper(key)
	if _, ok := u.urls[key]; !ok {
		panic("URL does not exist: " + key)
	}
	u.mu.Lock()
	delete(u.urls, key)
	for i, v := range u.order {
		if v == key {
			u.order = append(u.order[:i], u.order[i+1:]...)
			break
		}
	}
	u.mu.Unlock()
}

// Set display to none for all urls,
// or for a list of urls
func (u *URLs) Hide(urlname ...string) {
	if len(urlname) == 0 {
		for _, v := range u.urls {
			v.AttrStyle("display:none")
		}
		return
	}
	for _, v := range urlname {
		var url = u.Get(v)
		if url != nil {
			url.AttrStyle("display:none")
		}
	}
}

// Set display to param_display for all urls,
// or for a list of urls
func (u *URLs) Show(display string, urlname ...string) {
	if len(urlname) == 0 {
		for _, v := range u.urls {
			v.AttrStyle("display:" + display)
		}
		return
	}
	for _, v := range urlname {
		var url = u.Get(v)
		if url != nil {
			url.AttrStyle("display:" + display)
		}

	}
}

// Loop through all urls in order
func (u *URLs) InOrder(reverse ...bool) []*Element {
	var ret = make([]*Element, 0)
	if len(reverse) > 0 && reverse[0] {
		for i := len(u.order) - 1; i >= 0; i-- {
			ret = append(ret, u.urls[u.order[i]])
		}
		return ret
	}
	for _, v := range u.order {
		ret = append(ret, u.urls[v])
	}
	return ret
}

// Loop through all key, value urls in order
func (u *URLs) ForEach(f func(k string, elem *Element), reverse ...bool) {
	if len(reverse) > 0 && reverse[0] {
		for i := len(u.order) - 1; i >= 0; i-- {
			f(u.order[i], u.urls[u.order[i]])
		}
		return
	}
	for _, orderedKey := range u.order {
		f(orderedKey, u.urls[orderedKey])
	}
}

// Get the underlying map of URLs
func (u *URLs) Map() map[string]*Element {
	return u.urls
}

// Length of the underlying map of URLs
func (u *URLs) Len() int {
	return len(u.urls)
}

// Get the underlying slice of ordered keys
func (u *URLs) Keys() []string {
	return u.order
}

// Fill up the URLs map from a slice of Elements
func (u *URLs) FromElements(elems ...*Element) {
	for _, v := range elems {
		var href = v.GetAttr("href")
		if strings.HasPrefix(href, router.RT_PREFIX_EXTERNAL) {
			v.Delete("href")
			v.AttrHref(strings.TrimPrefix(href, router.RT_PREFIX_EXTERNAL))
			u.Set(v.Text, v, true)
			continue
		}
		u.Set(v.Text, v)
	}
}
