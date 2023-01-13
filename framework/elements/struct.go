//go:build js && wasm
// +build js,wasm

package elements

import (
	"strings"
	"sync"
	"syscall/js"

	"github.com/Nigel2392/jsext"
)

// https://github.com/golang/go/blob/master/src/html/escape.go
var htmlEscaper = strings.NewReplacer(
	`&`, "&amp;",
	`'`, "&#39;", // "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	`<`, "&lt;",
	`>`, "&gt;",
	`"`, "&#34;", // "&#34;" is shorter than "&quot;".
)

type Elements []*Element

// Element to be rendered in the DOM.
type Element struct {
	Tag                  string
	Text                 string
	textAfter            bool
	Attributes_Normal    map[string][]string
	Attributes_Boolean   map[string]bool
	Attributes_Semicolon map[string][]string
	Children             []*Element
	value                js.Value
	eventListeners       map[string][]func(this jsext.Value, event jsext.Event)
	animations           []Animation
}

func getText(t []string) string {
	if len(t) > 0 {
		return t[0]
	}
	return ""
}

// NewElement creates a new element.
// All normal HTML elements are predefined.
func NewElement(tag string, text ...string) *Element {
	return &Element{
		Tag:                  tag,
		Text:                 getText(text),
		Attributes_Normal:    make(map[string][]string),
		Attributes_Boolean:   make(map[string]bool),
		Attributes_Semicolon: make(map[string][]string),
		eventListeners:       make(map[string][]func(this jsext.Value, event jsext.Event)),
		animations:           make([]Animation, 0),
		Children:             make([]*Element, 0),
	}
}

// Return the js.Value of the element (if it has been rendered.)
func (e *Element) JSValue() js.Value {
	return e.value
}

// Return the jsext.Value of the element (if it has been rendered.)
func (e *Element) Value() jsext.Value {
	return jsext.Value(e.value)
}

// Return the jsext.Element of the element (if it has been rendered.)
func (e *Element) JSExtElement() jsext.Element {
	return jsext.Element(e.value)
}

// Report wether the element is of type specified.
func (e *Element) Is(tag string) bool {
	return strings.EqualFold(e.Tag, tag)
}

// Set the text to be rendered after all children are.
func (e *Element) TextAfter() *Element {
	e.textAfter = true
	return e
}

func (e *Element) Set(p string, v ...string) *Element {
	if e.value.IsUndefined() {
		var attrs = e.Attributes_Normal[p]
		attrs = append(attrs, v...)
		e.Attributes_Normal[p] = attrs
	} else {
		e.value.Call("setAttribute", p, strings.Join(v, " "))
	}
	return e
}

func (e *Element) SetBool(p string, v bool) *Element {
	if e.value.IsUndefined() {
		e.Attributes_Boolean[p] = v
	} else {
		e.value.Call("setAttribute", p, v)
	}
	return e
}

func (e *Element) SetSemicolon(p string, v ...string) *Element {
	if e.value.IsUndefined() {
		var attrs = e.Attributes_Semicolon[p]
		attrs = append(attrs, v...)
		e.Attributes_Semicolon[p] = attrs
	} else {
		e.value.Call("setAttribute", p, strings.Join(v, ";"))
	}
	return e
}

func (e *Element) Delete(p string) *Element {
	delete(e.Attributes_Normal, p)
	delete(e.Attributes_Boolean, p)
	delete(e.Attributes_Semicolon, p)
	if !e.value.IsUndefined() {
		e.value.Call("removeAttribute", p)
	}
	return e
}

func (e *Element) Add(p string, v ...string) *Element {
	if e.value.IsUndefined() {
		var attrs = e.Attributes_Normal[p]
		attrs = append(attrs, v...)
		e.Attributes_Normal[p] = attrs
		return e
	}
	var attr = e.value.Call("getAttribute", p)
	if attr.IsUndefined() {
		e.value.Call("setAttribute", p, strings.Join(v, " "))
	} else {

		var attrs = strings.Split(attr.String(), " ")
		var attrSlice = make([]string, 0)

		var attrMap = make(map[string]bool)
		for _, a := range attrs {
			attrMap[a] = true
		}
		for _, a := range v {
			attrMap[a] = true
		}
		for k := range attrMap {
			attrSlice = append(attrSlice, k)
		}

		e.value.Call("setAttribute", p, strings.Join(attrSlice, " "))
	}
	return e
}

func (e *Element) GetAttr(p string) string {
	if e.value.IsUndefined() {
		attr, ok := e.Attributes_Normal[p]
		if ok {
			return strings.Join(attr, " ")
		}
		attr, ok = e.Attributes_Semicolon[p]
		if ok {
			return strings.Join(attr, ";")
		}
	}
	var attr = e.value.Call("getAttribute", p)
	if attr.IsUndefined() {
		return ""
	}
	if attr.Type() == js.TypeBoolean {
		panic("attribute " + p + " is boolean")
	}
	return attr.String()
}

func (e *Element) GetBoolAttr(p string) bool {
	if e.value.IsUndefined() {
		var attr, ok = e.Attributes_Boolean[p]
		if ok {
			return attr
		}
	}
	var attr = e.value.Call("getAttribute", p)
	if attr.IsUndefined() {
		return false
	}
	if attr.Type() != js.TypeBoolean {
		panic("attribute " + p + " is not boolean")
	}
	return attr.Bool()
}

func (e *Element) Append(children ...*Element) *Element {
	if e.value.IsUndefined() || e.value.IsNull() {
		e.Children = append(e.Children, children...)
	} else {
		for _, child := range children {
			e.JSExtElement().AppendChild(child.Render())
		}
	}
	return e
}

func (e *Element) Prepend(children ...*Element) *Element {
	if e.value.IsUndefined() {
		e.Children = append(children, e.Children...)
	} else {
		for _, child := range children {
			e.value.Call("prepend", child.Render())
		}
	}
	return e
}

func (e *Element) AppendAfter(id string, children ...*Element) *Element {
	if e.value.IsUndefined() {
		for i, child := range e.Children {
			if child.GetAttr("id") == id {
				e.Children = append(e.Children[:i+1], append(children, e.Children[i+1:]...)...)
				return e
			}
		}
	} else {
		var after = e.value.Call("querySelector", "#"+id)
		if after.IsUndefined() {
			return e
		}
		for _, child := range children {
			after.Call("after", child.Render())
		}
	}
	return e
}

func (e *Element) AppendBefore(id string, children ...*Element) *Element {
	if e.value.IsUndefined() {
		for i, child := range e.Children {
			if child.GetAttr("id") == id {
				e.Children = append(e.Children[:i], append(children, e.Children[i:]...)...)
				return e
			}
		}
	} else {
		var before = e.value.Call("querySelector", "#"+id)
		if before.IsUndefined() {
			return e
		}
		for _, child := range children {
			before.Call("before", child.Render())
		}
	}
	return e
}

func (e *Element) AppendAfterElement(element *Element, children ...*Element) *Element {
	if e.value.IsUndefined() || element.value.IsUndefined() {
		for i, child := range e.Children {
			if child == element {
				e.Children = append(e.Children[:i+1], append(children, e.Children[i+1:]...)...)
				return e
			}
		}
	} else {
		for _, child := range children {
			element.value.Call("after", child.Render())
		}
	}
	return e
}

func (e *Element) AppendBeforeElement(element *Element, children ...*Element) *Element {
	if e.value.IsUndefined() {
		for i, child := range e.Children {
			if child == element {
				e.Children = append(e.Children[:i], append(children, e.Children[i:]...)...)
				return e
			}
		}
	} else {
		for _, child := range children {
			element.value.Call("before", child.Render())
		}
	}
	return e
}

// Loop over all inner elements recursively asynchronously.
func (e *Element) AsyncForEach(fn func(*Element) bool) {
	var wg sync.WaitGroup
	var terminate = make(chan struct{}, 1)
	wg.Add(1)
	go e.asyncForEach(fn, &wg, terminate)
	wg.Wait()
	for len(terminate) > 0 {
		<-terminate
	}
	close(terminate)
}

func (e *Element) asyncForEach(fn func(*Element) bool, wg *sync.WaitGroup, terminate chan struct{}) {
	defer wg.Done()
	if fn(e) {
		terminate <- struct{}{}
		return
	}
	select {
	case <-terminate:
		return
	default:
		wg.Add(len(e.Children))
		for _, child := range e.Children {
			go child.asyncForEach(fn, wg, terminate)
		}
	}
}

func (e *Element) AddEventListener(event string, fn func(this jsext.Value, event jsext.Event)) {
	if e.value.IsUndefined() {
		e.eventListeners[event] = append(e.eventListeners[event], fn)
		return
	}
	e.JSExtElement().AddEventListener(event, fn)
}

// If the element is not rendered yet, we will delete the event key from the eventListeners map.
// This means that all functions will be removed from the event.
// This is due to the fact that we cannot compare functions in Go.
//   - Will panic if the element is already rendered and fn is nil.
//   - Will not delete the eventlistener from the element's eventListener map once it is rendered.
func (e *Element) RemoveEventListener(event string, fn func(this jsext.Value, event jsext.Event)) {
	if e.value.IsUndefined() || e.value.IsNull() {
		delete(e.eventListeners, event)
		return
	}

	if fn == nil {
		panic("fn cannot be nil")
	}

	e.JSExtElement().RemoveEventListener(event, fn)
}

// Set the text of the element
// Autoescape possible HTML input.
func (e *Element) InnerText(text string) *Element {
	if e.value.IsUndefined() {
		e.Text = htmlEscaper.Replace(text)
	} else {
		e.value.Set("innerHTML", htmlEscaper.Replace(text))
	}
	return e
}

// Set raw HTML of the element.
func (e *Element) InnerHTML(html string) *Element {
	if e.value.IsUndefined() {
		e.Text = html
	} else {
		e.value.Set("innerHTML", html)
	}
	return e
}

func (e *Element) Remove() {
	if e.value.IsUndefined() {
		return
	}
	e.value.Call("remove")
}

func (e *Element) ClearInnerHTML() {
	if e.value.IsUndefined() {
		return
	}
	e.value.Set("innerHTML", "")
}

func (e *Element) ClearInnerText() {
	if e.value.IsUndefined() {
		return
	}
	e.value.Set("innerText", "")
}

func (e *Element) Clear() {
	if e.value.IsUndefined() {
		return
	}
	e.value.Set("innerHTML", "")
	e.value.Set("innerText", "")
	for _, child := range e.Children {
		child.Remove()
	}
}

func (e *Element) Render() jsext.Element {
	if e.value.IsUndefined() {
		e.generate(js.Undefined())
	}
	return e.JSExtElement()
}

func (e *Element) RenderTo(appendToQuerySelector ...string) *Element {
	// Create the element
	if len(appendToQuerySelector) > 0 {
		var parent = js.Global().Get("document").Call("querySelector", appendToQuerySelector[0])
		if !parent.IsUndefined() && !parent.IsNull() {
			e.generate(parent)
			return e
		}
	}
	e.generate(js.Undefined())
	return e
}

func (e *Element) generate(parent js.Value) js.Value {
	// Create the element
	e.value = js.Global().Get("document").Call("createElement", e.Tag)

	// Append to parent
	if !parent.IsUndefined() && !parent.IsNull() {
		parent.Call("appendChild", e.value)
	}

	// Set attributes
	for k, attr := range e.Attributes_Normal {
		e.value.Call("setAttribute", k, strings.Join(attr, " "))
	}
	// Set attributes delimited by semicolon
	for k, attr := range e.Attributes_Semicolon {
		e.value.Call("setAttribute", k, strings.Join(attr, ";"))
	}
	// Set boolean attributes
	for k, attr := range e.Attributes_Boolean {
		e.value.Call("setAttribute", k, attr)
	}
	//  Set text before children
	if !e.textAfter {
		if e.Text != "" {
			e.value.Call("insertAdjacentHTML", "afterbegin", e.Text)
		}
	}
	// Generate children
	for _, child := range e.Children {
		child.generate(e.value)
	}
	// Set text after children
	if e.textAfter {
		if e.Text != "" {
			e.value.Call("insertAdjacentHTML", "beforeend", e.Text)
		}
	}

	// Add event listeners
	for event, listeners := range e.eventListeners {
		var elem = jsext.Element(e.value)
		for _, listener := range listeners {
			elem.AddEventListener(event, listener)
		}
	}

	e.animate()

	return e.value
}

//func (e *Element) asyncSetAttrs() {
//	// Create a channel for a maximum of 10 elements
//	var guard = make(chan struct{}, 10)
//	// Create a waitgroup
//	var wg sync.WaitGroup
//	wg.Add(3)
//	// Loop over all attributes asynchronously
//	go func() {
//		defer wg.Done()
//		for k, attr := range e.Attributes_Normal {
//			// Add to the channel
//			guard <- struct{}{}
//			// Create a new goroutine
//			go func(key string, attribute []string) {
//				// Set the attribute
//				e.value.Call("setAttribute", key, strings.Join(attribute, " "))
//				<-guard
//			}(k, attr)
//		}
//		// Wait for all goroutines to finish
//	}()
//	go func() {
//		defer wg.Done()
//		for k, attr := range e.Attributes_Semicolon {
//			guard <- struct{}{}
//			go func(key string, attribute []string) {
//				e.value.Call("setAttribute", key, strings.Join(attribute, ";"))
//				<-guard
//			}(k, attr)
//		}
//	}()
//	go func() {
//		defer wg.Done()
//		for k, attr := range e.Attributes_Boolean {
//			guard <- struct{}{}
//			go func(key string, attribute bool) {
//				e.value.Call("setAttribute", key, attribute)
//				<-guard
//			}(k, attr)
//		}
//	}()
//	// Wait for all goroutines to finish
//	wg.Wait()
//	close(guard)
//}

//type Attribute struct {
//	value any
//}
//
//func (a *Attribute) String() string {
//	return a.value.(string)
//}
//
//func (a *Attribute) Bool() bool {
//	return a.value.(bool)
//}
//
//func (a *Attribute) Int() int {
//	return a.value.(int)
//}
//
//func (a *Attribute) Int64() int64 {
//	return a.value.(int64)
//}
//
//func (a *Attribute) Float64() float64 {
//	return a.value.(float64)
//}
//
//func (e *Element) GetAttr(p string) Attribute {
//	if e.value.IsUndefined() {
//		attr, ok := e.Attributes_Normal[p]
//		if ok {
//			return Attribute{strings.Join(attr, " ")}
//		}
//		attr, ok = e.Attributes_Semicolon[p]
//		if ok {
//			return Attribute{strings.Join(attr, ";")}
//		}
//	}
//	var attr = e.value.Call("getAttribute", p)
//	if attr.IsUndefined() {
//		return ""
//	}
//	if attr.Type() == js.TypeBoolean {
//		return Attribute{attr.Bool()}
//	}
//	return Attribute{attr.String()}
//}
//
