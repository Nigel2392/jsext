//go:build js && wasm
// +build js,wasm

package elements

import (
	"strconv"
	"strings"
	"sync"
	"syscall/js"

	"github.com/Nigel2392/jsext"
)

const SPACING string = "    "

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
	// Element tag.
	Tag string

	// Text to be rendered in the element.
	Text string

	// Render the text after the children.
	textAfter bool

	// Element attributes.
	Attributes_Normal    map[string][]string
	Attributes_Boolean   map[string]bool
	Attributes_Semicolon map[string][]string

	// Children of the element.
	Children []*Element

	// Underlying javascript element.
	value js.Value

	// Eventlisteners for the element.
	// These are only filled when the element is not rendered.
	eventListeners map[string][]func(this jsext.Value, event jsext.Event)

	// Element animations.
	// These are only filled when the element is not rendered.
	Animations *Animations

	// Function to call when the javascript element is ready.
	onRender []func(this jsext.Element)

	// Wether to render the end tag or not.
	// This is only used for rendering the element as a string.
	noEnd bool
}

func getText(t []string) string {
	if len(t) > 0 {
		return t[0]
	}
	return ""
}

// Do not render the end tag.
func (e *Element) NoEnd() {
	e.noEnd = true
}

func (e *Element) OnRender(f func(this jsext.Element)) {
	e.onRender = append(e.onRender, f)
}

// NewElement creates a new element.
// All normal HTML elements are predefined.
func NewElement(tag string, text ...string) *Element {
	var e = &Element{
		Tag:                  tag,
		Text:                 getText(text),
		Attributes_Normal:    make(map[string][]string),
		Attributes_Boolean:   make(map[string]bool),
		Attributes_Semicolon: make(map[string][]string),
		eventListeners:       make(map[string][]func(this jsext.Value, event jsext.Event)),
		Children:             make([]*Element, 0),
		onRender:             make([]func(this jsext.Element), 0),
	}
	e.Animations = &Animations{element: e, animations: make([]Animation, 0)}
	return e
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

// Set an attribute on the element, overwriting any previous value.
func (e *Element) SetAttr(p string, v ...string) *Element {
	if e.value.IsUndefined() {
		e.Attributes_Normal[p] = v
	} else {
		e.value.Call("setAttribute", p, strings.Join(v, " "))
	}
	return e
}

// Set an attribute of the element.
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

// Set a boolean attribute on the element.
func (e *Element) SetBool(p string, v bool) *Element {
	if e.value.IsUndefined() {
		e.Attributes_Boolean[p] = v
	} else {
		e.value.Call("setAttribute", p, v)
	}
	return e
}

// Add a semicolon separated attribute on the element.
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

// Delete an attribute from the element.
func (e *Element) Delete(p string) *Element {
	delete(e.Attributes_Normal, p)
	delete(e.Attributes_Boolean, p)
	delete(e.Attributes_Semicolon, p)
	if !e.value.IsUndefined() {
		e.value.Call("removeAttribute", p)
	}
	return e
}

// Add an attribute to the element.
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

// Get an attribute of the element.
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

// Get a boolean attribute of the element.
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

// Append appends children to the elements children.
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

// Prepend prepends children to the elements children.
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

// AppendAfter appends children after the given element by id.
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

// AppendBefore appends children before the given element by id.
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

// AppendAfterElement appends children after the given element.
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

// AppendBeforeElement appends children before the given element.
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

// Loop over all inner elements recursively asynchronously.
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

// Add an eventlistener to the element.
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

// Remove the element from the javascript dom.
func (e *Element) Remove() {
	if !e.value.IsUndefined() {
		e.value.Call("remove")
	}
	e = nil
}

// Clear the inner html of the element.
// This includes removing all element children.
func (e *Element) ClearInnerHTML() {
	e.Text = ""
	e.Children = make([]*Element, 0)
	if !e.value.IsUndefined() {
		e.value.Set("innerHTML", "")
	}
}

// Clear the inner text of the element.
func (e *Element) ClearInnerText() {
	e.Text = ""
	if !e.value.IsUndefined() {
		e.value.Set("innerText", "")
	}
}

// Reset all element attributes or
// clear the javascript element and remove all children.
func (e *Element) Reset() {
	// Remove all attributes
	e.Attributes_Normal = make(map[string][]string, 0)
	e.Attributes_Semicolon = make(map[string][]string, 0)
	e.Attributes_Boolean = make(map[string]bool, 0)
	// Clear the inner html
	e.Text = ""
	e.Children = make([]*Element, 0)
	// Remove all animations
	e.Animations.animations = make([]Animation, 0)
	// Remove all event listeners
	e.eventListeners = make(map[string][]func(this jsext.Value, event jsext.Event), 0)
	if !e.value.IsUndefined() {
		// Clear the inner html
		e.value.Set("innerHTML", "")
		e.value.Set("innerText", "")
		// Remove all attributes
		elem := e.JSExtElement()
		attrs := elem.Get("attributes")
		for i := 0; i < attrs.Length(); i++ {
			elem.Call("removeAttribute", attrs.Index(i).Get("name").Value())
		}

		// Remove all animations
		var animations = jsext.Call("getAnimations", e.value)
		for i := 0; i < animations.Length(); i++ {
			animations.Index(i).Call("cancel")
		}

		// Remove all event listeners
		var eventListeners = jsext.Call("getEventListeners", e.value)
		var keys = jsext.Call("Object", "keys", eventListeners.Value())
		for i := 0; i < keys.Length(); i++ {
			var key = keys.Index(i).String()
			var eventListener = eventListeners.Get(key)
			for j := 0; j < eventListener.Length(); j++ {
				e.value.Call("removeEventListener", key, eventListener.Index(j).Get("listener").Value())
			}
		}
	}
}

// Generate the javascript element and return it.
func (e *Element) Render() jsext.Element {
	if e.value.IsUndefined() {
		e.generate(js.Undefined())
	}
	return e.JSExtElement()
}

// Render the element onto another javascript element by using a query selector.
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

// Generate the javascript element and return it.
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

	// OnRender
	if len(e.onRender) > 0 {
		for _, fn := range e.onRender {
			fn(jsext.Element(e.value))
		}
	}

	// Animate
	e.Animations.animate()

	return e.value
}

// Buffer interface.
// Used to generate the HTML.
type Buffer interface {
	WriteString(string) (int, error)
	String() string
	Reset()
	Len() int
	Write(p []byte) (int, error)
	Grow(n int)
}

// Generate the Element and children as HTML.
func (e *Element) Generate(spacing string, buf Buffer) {
	var length int
	if e.Tag != "" {
		length += len(e.Tag) + len(spacing) + 4
		if e.Text != "" {
			length += len(e.Text) + len(spacing) + len(SPACING) + 2
		}
		if !e.noEnd {
			length += len(spacing) + 5 + len(e.Tag)
		}
		buf.Grow(length)

		buf.WriteString(spacing)
		buf.WriteString("<")
		buf.WriteString(e.Tag)

		writeAttrs(buf, e.Attributes_Normal, " ")
		writeAttrs(buf, e.Attributes_Semicolon, ";")

		for k, attr := range e.Attributes_Boolean {
			// buf.WriteString(" " + k + "=\"" + strconv.FormatBool(attr) + "\"")
			if attr {
				buf.Grow(10 + len(k))
			} else {
				buf.Grow(11 + len(k))
			}
			buf.WriteString(" ")
			buf.WriteString(k)
			buf.WriteString("=\"")
			buf.WriteString(strconv.FormatBool(attr))
			buf.WriteString("\"")
		}
		buf.WriteString(">\n")
	}

	if !e.textAfter {
		if e.Text != "" {
			// buf.WriteString(spacing + SPACING + e.text + "\n")
			buf.WriteString(spacing)
			buf.WriteString(SPACING)
			buf.WriteString(e.Text)
			buf.WriteString("\n")
		}
	}
	// Loop over inner html elements
	for _, e := range e.Children {
		e.Generate(spacing+SPACING, buf)
	}

	if e.textAfter {
		if e.Text != "" {
			// buf.WriteString(spacing + SPACING + e.text + "\n")
			buf.WriteString(spacing)
			buf.WriteString(SPACING)
			buf.WriteString(e.Text)
			buf.WriteString("\n")
		}
	}
	if e.Tag != "" {
		if e.noEnd {
			// buf.WriteString(spacing + "</" + e.Type + ">\n")
			buf.WriteString(spacing)
			buf.WriteString("</")
			buf.WriteString(e.Tag)
			buf.WriteString(">\n")
		}
	}
}

// Write the attributes to the buffer.
func writeAttrs(buf Buffer, attrs map[string][]string, sep string) {
	for k, attr := range attrs {
		// buf.WriteString(" " + k + "=\"" + strings.Join(attr, sep) + "\"")
		n := len(sep)*(len(attr)-1) + len(k) + 6
		for i := 0; i < len(attr); i++ {
			n += len(attr[i])
		}
		buf.Grow(n)
		buf.WriteString(" ")
		buf.WriteString(k)
		buf.WriteString("=\"")
		// buf.WriteString(strings.Join(attr, sep))
		if len(attr) == 1 {
			buf.WriteString(attr[0])
		} else {
			buf.WriteString(attr[0])
			for _, s := range attr[1:] {
				buf.WriteString(sep)
				buf.WriteString(s)
			}
		}
		buf.WriteString("\"")
	}
}
