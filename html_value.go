//go:build js && wasm
// +build js,wasm

package jsext

import "syscall/js"

// Defines multiple elements.
// This is a wrapper around a slice of []Element(js.Value).
type Elements []Element

// Get the length of the slice.
func (e Elements) Len() int {
	return len(e)
}

// Return all inner js.Value in the slice.
func (e Elements) Value() []js.Value {
	var values []js.Value = make([]js.Value, len(e))
	for i, v := range e {
		values[i] = v.JSValue()
	}
	return values
}

// AddEventListener adds an event listener to all elements in the slice.
func (e Elements) AddEventListener(event string, listener func(this Value, event Event)) Elements {
	for _, el := range e {
		el.AddEventListener(event, listener)
	}
	return e
}

// RemoveEventListener removes an event listener from all elements in the slice.
func (e Elements) RemoveEventListener(event string, listener func(this Value, event Event)) Elements {
	for _, el := range e {
		el.RemoveEventListener(event, listener)
	}
	return e
}

// Element is a wrapper around js.Value.
type Element js.Value

func (e Element) IsNull() bool {
	return js.Value(e).IsNull()
}

func (e Element) IsUndefined() bool {
	return js.Value(e).IsUndefined()
}

func (e Element) IsNaN() bool {
	return js.Value(e).IsNaN()
}

// JSValue returns the underlying js.Value.
func (e Element) JSValue() js.Value {
	return js.Value(e)
}

// Value returns as a Value(js.Value) wrapper.
func (e Element) Value() Value {
	return Value(e.JSValue())
}

// Set sets a property on the element.
func (e Element) Set(p string, v interface{}) Element {
	e.JSValue().Set(p, v)
	return e
}

// Get gets a property from the element.
func (e Element) Get(p string) Value {
	return Value(e.JSValue().Get(p))
}

// Call calls a method on the element.
func (e Element) Call(m string, args ...interface{}) Value {
	return Value(e.JSValue().Call(m, args...))
}

// CallFunc is used by state management.
func (e Element) CallFunc(name string, args ...interface{}) {
	e.Call(name, args...)
}

// Delete deletes a property from the element.
func (e Element) Delete(p string) {
	e.JSValue().Delete(p)
}

// Equal returns true if the element is equal to the other js.Value.
func (e Element) Equal(other js.Value) bool {
	return e.JSValue().Equal(other)
}

// AddEventListener adds an event listener to the element.
func (e Element) AddEventListener(event string, listener func(this Value, event Event)) Element {
	e.JSValue().Call("addEventListener", event, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		listener(Value(this), Event(args[0]))
		return nil
	}))
	return e
}

// RemoveEventListener removes an event listener from the element.
func (e Element) RemoveEventListener(event string, listener func(this Value, event Event)) Element {
	e.JSValue().Call("removeEventListener", event, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		listener(Value(this), Event(args[0]))
		return nil
	}))
	return e
}

// AppendChild appends a child to the element.
func (e Element) AppendChild(child Element) Element {
	e.JSValue().Call("appendChild", child.JSValue())
	return e
}

// AppendChildren appends multiple children to the element.
func (e Element) AppendChildren(children []Element) Element {
	for _, child := range children {
		e.AppendChild(child)
	}
	return e
}

// RemoveChild removes a child from the element.
func (e Element) RemoveChild(child Element) Element {
	e.JSValue().Call("removeChild", child.JSValue())
	return e
}

// RemoveChildren removes multiple children from the element.
func (e Element) RemoveChildren(children []Element) Element {
	for _, child := range children {
		e.RemoveChild(child)
	}
	return e
}

// Remove removes the element from the DOM.
func (e Element) SetAttribute(name, value string) Element {
	e.JSValue().Call("setAttribute", name, value)
	return e
}

// InnerHTML sets the inner HTML of the element.
func (e Element) InnerHTML(html string) Element {
	e.JSValue().Set("innerHTML", html)
	return e
}

// InnerText sets the inner text of the element.
func (e Element) InnerText(text string) Element {
	e.JSValue().Set("innerText", text)
	return e
}

// SetInnerElement clears the inner HTML and appends the element.
func (e Element) InnerElement(el Element) Element {
	e.JSValue().Set("innerHTML", "")
	e.AppendChild(el)
	return e
}

// Get the style of the element.
func (e Element) Style() Style {
	return Style(e.Get("style"))
}

// Set the style of an element property.
func (e Element) StyleProperty(name, value string) Element {
	e.JSValue().Get("style").Set(name, value)
	return e
}

// Get the style of an element property.
func (e Element) GetStyleProperty(name string) string {
	return e.JSValue().Get("style").Get(name).String()
}

// Get the value of the element.
func (e Element) GetClassList() Value {
	return Value(e.JSValue().Get("classList"))
}

// Set multiple classes on the element.
func (e Element) ClassList(c ...string) Value {
	var cList = e.Get("classList")
	if len(c) == 0 {
		return cList
	}
	for _, cl := range c {
		cList.Call("add", cl)
	}
	return cList
}

// Get the parentElement
func (e Element) ParentElement() Element {
	return Element(e.JSValue().Get("parentElement"))
}

// Get the children of the element.
func (e Element) GetChildren() Value {
	return Value(e.JSValue().Get("children"))
}

// Get the children of the element.
func (e Element) GetChildNodes() Value {
	return Value(e.JSValue().Get("childNodes"))
}

// Get the ID of the element.
func (e Element) ID() string {
	return e.JSValue().Get("id").String()
}

// Set the ID of the element.
func (e Element) SetID(id string) Element {
	e.JSValue().Set("id", id)
	return e
}

// Get the className of the element.
func (e Element) ClassName() string {
	return e.JSValue().Get("className").String()
}

// Set the className of the element.
func (e Element) SetClassName(className string) Element {
	e.JSValue().Set("className", className)
	return e
}

// Get an inner element by ID.
func (e Element) GetElementById(id string) Element {
	return Element(e.JSValue().Call("getElementById", id))
}

// Get an inner element by class name.
func (e Element) GetElementsByClassName(className string) Element {
	return Element(e.Call("getElementsByClassName", className))
}

// Get an inner element by tag name.
func (e Element) GetElementsByTagName(tagName string) Element {
	return Element(e.Call("getElementsByTagName", tagName))
}

// Get the scroll height of the element.
func (e Element) ScrollHeight() int {
	return e.JSValue().Get("scrollHeight").Int()
}

// Get the scroll width of the element.
func (e Element) ScrollWidth() int {
	return e.JSValue().Get("scrollWidth").Int()
}

// Get the scroll top of the element.
func (e Element) ScrollTop() int {
	return e.JSValue().Get("scrollTop").Int()
}

// Get the scroll left of the element.
func (e Element) ScrollLeft() int {
	return e.JSValue().Get("scrollLeft").Int()
}

// Set the sccrollTo of the element.
func (e Element) ScrollTo(x, y int) {
	e.JSValue().Call("scrollTo", x, y)
}

// Scroll the element into view.
func (e Element) ScrollIntoView(center bool) {
	e.JSValue().Call("scrollIntoView", center)
}

// Scroll the element into view if needed.
func (e Element) ScrollIntoViewIfNeeded(center bool) {
	e.JSValue().Call("scrollIntoViewIfNeeded", center)
}

// Get the clientWidth of the element.
func (e Element) Width() int {
	return e.JSValue().Get("clientWidth").Int()
}

// Get the clientHeight of the element.
func (e Element) Height() int {
	return e.JSValue().Get("clientHeight").Int()
}

// Get the assigned slot of the element.
func (e Element) AssignedSlot() Element {
	return Element(e.JSValue().Get("assignedSlot"))
}

// Get attributes of the element.
func (e Element) Attributes() Value {
	return Value(e.JSValue().Get("attributes"))
}

// Get the childElementCount of the element.
func (e Element) ChildElementCount() int {
	return e.JSValue().Get("childElementCount").Int()
}

// Get the children of the element.
func (e Element) Children() Value {
	return Value(e.JSValue().Get("children"))
}

// Get the clientHeight of the element.
func (e Element) ClientHeight() int {
	return e.JSValue().Get("clientHeight").Int()
}

// Get the clientLeft of the element.
func (e Element) ClientLeft() int {
	return e.JSValue().Get("clientLeft").Int()
}

// Get the clientTop of the element.
func (e Element) ClientTop() int {
	return e.JSValue().Get("clientTop").Int()
}

// Get the clientWidth of the element.
func (e Element) ClientWidth() int {
	return e.JSValue().Get("clientWidth").Int()
}

// Get the element's dataset.
func (e Element) Dataset() Value {
	return Value(e.JSValue().Get("dataset"))
}

// Return the element's dataset as map
func (e Element) MapDataset() map[string]string {
	var dataset = e.Dataset()
	return ObjectToMapT[string](dataset.Value())
}

// Get the element's first child.
func (e Element) FirstElementChild() Element {
	return Element(e.JSValue().Get("firstElementChild"))
}

// Get the element's href.
func (e Element) Href() string {
	return e.JSValue().Get("href").String()
}

// Set elements after the element.
func (e Element) After(elements ...Element) {
	for _, element := range elements {
		e.JSValue().Call("after", element.JSValue())
	}
}

// Set elements before the element.
func (e Element) Before(elements ...Element) {
	for _, element := range elements {
		e.JSValue().Call("before", element.JSValue())
	}
}

// Append elements to the element.
func (e Element) Append(elements ...Element) {
	for _, element := range elements {
		e.JSValue().Call("append", element.JSValue())
	}
}

// Append the element to the parent.
func (e Element) AppendTo(parent Element) {
	parent.AppendChild(e)
}

// Append the element to the parent.
func (e Element) Prepend(elements ...Element) {
	for _, element := range elements {
		e.JSValue().Call("prepend", element.JSValue())
	}
}

// Insert the element before the before element.
func (e Element) InsertBefore(element, before Element) {
	e.JSValue().Call("insertBefore", element.JSValue(), before.JSValue())
}

// Replace the element with the before element.
func (e Element) ReplaceChild(element, before Element) {
	e.JSValue().Call("replaceChild", element.JSValue(), before.JSValue())
}

// Remove the element.
func (e Element) Remove() {
	e.JSValue().Call("remove")
}

// Animate the element.
func (e Element) Animate(keyframes []interface{}, options map[string]interface{}) Value {
	return Value(e.JSValue().Call("animate", SliceToArray(keyframes).Value(), MapToObject(options).Value()))
}
