package jse

import (
	"strings"
	"syscall/js"

	"github.com/Nigel2392/jsext/v2"
)

type Element jsext.Element

func NewElement(tag string, text ...string) *Element {
	var e = jsext.CreateElement(tag)
	if len(text) > 0 {
		e.InnerHTML(strings.Join(text, "\n"))
	}
	return (*Element)(&e)
}

type JavascriptConstraint interface {
	jsext.Element | jsext.Value | SVG | FormElement | SelectElement | Element | *Element | *SVG | *FormElement | *SelectElement
}

func Make[T JavascriptConstraint](e T) *Element {
	switch e := any(e).(type) {
	case jsext.Element:
		return (*Element)(&e)
	case jsext.Value:
		return (*Element)(&e)
	case SVG:
		return (*Element)(&e)
	case FormElement:
		return (*Element)(&e)
	case SelectElement:
		return (*Element)(&e)
	case Element:
		return &e
	case *Element:
		return e
	case *SVG:
		return (*Element)(e)
	case *FormElement:
		return (*Element)(e)
	case *SelectElement:
		return (*Element)(e)
	}
	panic("unreachable")
}

func (e *Element) Element() jsext.Element {
	return (jsext.Element)(*e)
}

func (e *Element) Value() jsext.Value {
	return (jsext.Value)(*e)
}

func (e *Element) JSValue() js.Value {
	return (js.Value)(*e)
}

func (e *Element) AppendChild(children ...*Element) *Element {
	for _, child := range children {
		e.Call("appendChild", child.JSValue())
	}
	return e
}

func (e *Element) PrependChild(children ...*Element) *Element {
	for _, child := range children {
		e.Call("insertBefore", child.JSValue(), e.JSValue().Get("firstChild"))
	}
	return e
}

func (e *Element) NewElement(typ string, innerText ...string) *Element {
	var elem = NewElement(typ, innerText...)
	e.AppendChild(elem)
	return elem
}

func (e *Element) ClearInnerHTML() {
	e.JSValue().Set("innerHTML", "")
}

func (e *Element) InnerHTML(s string) *Element {
	e.Element().InnerHTML(s)
	return e
}

func (e *Element) InnerText(s string) *Element {
	e.Element().InnerText(s)
	return e
}

func (e *Element) InlineClasses(classes ...string) *Element {
	e.Element().ClassList(classes...)
	return e
}

func (e *Element) ClassList(s ...string) jsext.Value {
	return e.Element().ClassList(s...)
}

func (e *Element) Style() jsext.Style {
	return e.Element().Style()
}

func (e *Element) SetAttrMap(m map[string]string) *Element {
	for k, v := range m {
		e.SetAttr(k, v)
	}
	return e
}

func (e *Element) SetAttr(p string, s ...string) *Element {
	e.Value().Call("setAttribute", p, strings.Join(s, " "))
	return e
}

func (e *Element) GetAttr(p string) string {
	return e.Value().Call("getAttribute", p).String()
}

func (e *Element) DelAttr(p string) *Element {
	e.Value().Call("removeAttribute", p)
	return e
}

// RemoveChild removes a child from the Element
func (e *Element) RemoveChild(child Element) *Element {
	e.JSValue().Call("removeChild", child.JSValue())
	return e
}

// RemoveChildren removes multiple children from the Element
func (e *Element) RemoveChildren(children []Element) *Element {
	for _, child := range children {
		e.RemoveChild(child)
	}
	return e
}

// Get the parentElement
func (e *Element) ParentElement() *Element {
	var p = e.JSValue().Get("parentElement")
	return (*Element)(&p)
}

// Get an inner element by ID.
func (e *Element) GetElementById(id string) *Element {
	var elem = e.JSValue().Call("getElementById", id)
	return (*Element)(&elem)
}

// Get an inner element by class name.
func (e *Element) GetElementsByClassName(className string) *Element {
	var elem = e.Call("getElementsByClassName", className)
	return (*Element)(&elem)
}

// Get an inner element by tag name.
func (e *Element) GetElementsByTagName(tagName string) *Element {
	var elem = e.Call("getElementsByTagName", tagName)
	return (*Element)(&elem)
}

// Add an event listener to the Element
//
// This will return the function that was added to the element.
func (e *Element) AddEventListener(event string, callback func(this *Element, event jsext.Event)) js.Func {
	if e == nil {
		return js.Func{Value: js.Null()}
	}
	if callback == nil {
		return js.Func{Value: js.Null()}
	}

	var f = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return nil
		}
		callback(e, jsext.Event(args[0]))
		return nil
	})

	e.JSValue().Call("addEventListener", event, f)

	return f
}

// Get the scroll height of the Element
func (e *Element) ScrollHeight() int {
	return e.JSValue().Get("scrollHeight").Int()
}

// Get the scroll width of the Element
func (e *Element) ScrollWidth() int {
	return e.JSValue().Get("scrollWidth").Int()
}

// Get the scroll top of the Element
func (e *Element) ScrollTop() int {
	return e.JSValue().Get("scrollTop").Int()
}

// Get the scroll left of the Element
func (e *Element) ScrollLeft() int {
	return e.JSValue().Get("scrollLeft").Int()
}

// Set the sccrollTo of the Element
func (e *Element) ScrollTo(x, y int) {
	e.JSValue().Call("scrollTo", x, y)
}

// Scroll the *Element into view.
func (e *Element) ScrollIntoView(center bool) {
	e.JSValue().Call("scrollIntoView", center)
}

// Scroll the *Element into view if needed.
func (e *Element) ScrollIntoViewIfNeeded(center bool) {
	e.JSValue().Call("scrollIntoViewIfNeeded", center)
}

// Get the clientWidth of the Element
func (e *Element) Width() int {
	return e.JSValue().Get("clientWidth").Int()
}

// Get the clientHeight of the Element
func (e *Element) Height() int {
	return e.JSValue().Get("clientHeight").Int()
}

// Get the children of the Element
func (e *Element) Children() []*Element {
	var children = e.JSValue().Get("children")
	var length = children.Length()
	var elems = make([]*Element, length)
	for i := 0; i < length; i++ {
		var child = children.Index(i)
		elems[i] = (*Element)(&child)
	}
	return elems
}

// Get the clientHeight of the Element
func (e *Element) ClientHeight() int {
	return e.JSValue().Get("clientHeight").Int()
}

// Get the clientLeft of the Element
func (e *Element) ClientLeft() int {
	return e.JSValue().Get("clientLeft").Int()
}

// Get the clientTop of the Element
func (e *Element) ClientTop() int {
	return e.JSValue().Get("clientTop").Int()
}

// Get the clientWidth of the Element
func (e *Element) ClientWidth() int {
	return e.JSValue().Get("clientWidth").Int()
}

// Get the *Element's dataset.
func (e *Element) Dataset() jsext.Value {
	return jsext.Value(e.JSValue().Get("dataset"))
}

// Return the *Element's dataset as map
func (e *Element) MapDataset() map[string]string {
	var dataset = e.Dataset()
	return jsext.ObjectToMapT[string](dataset.Value())
}

// Get the *Element's first child.
func (e *Element) FirstElementChild() Element {
	return Element(e.JSValue().Get("firstElementChild"))
}

// Insert the *Element before the before Element
func (e *Element) InsertBefore(element, before *Element) {
	e.JSValue().Call("insertBefore", element.JSValue(), before.JSValue())
}

// Replace the *Element with the before Element
func (e *Element) ReplaceChild(element, before *Element) {
	e.JSValue().Call("replaceChild", element.JSValue(), before.JSValue())
}

// Remove the Element
func (e *Element) Remove() {
	e.JSValue().Call("remove")
}

// Animate the Element
func (e *Element) Animate(keyframes []interface{}, options map[string]interface{}) jsext.Value {
	return jsext.Value(e.JSValue().Call("animate", jsext.SliceToArray(keyframes).Value(), jsext.MapToObject(options).Value()))
}

func (e *Element) SetMap(m map[string]string) *Element {
	for k, v := range m {
		e.Set(k, v)
	}
	return e
}

// /////////////////////////////////////////////////////////
//
// js.Value methods.
//
// /////////////////////////////////////////////////////////
func (e *Element) Bool() bool {
	return e.Value().Bool()
}
func (e *Element) Call(m string, args ...any) jsext.Value {
	return jsext.Value(e.Value().Call(m, args...))
}
func (e *Element) Delete(p string) {
	e.Value().Delete(p)
}
func (e *Element) Equal(other js.Value) bool {
	return e.Value().Equal(other)
}
func (e *Element) Float() float64 {
	return e.Value().Float()
}
func (e *Element) Get(p string) jsext.Value {
	return jsext.Value(e.Value().Get(p))
}
func (e *Element) Index(i int) jsext.Value {
	return jsext.Value(e.Value().Index(i))
}
func (e *Element) Int() int {
	return e.Value().Int()
}
func (e *Element) Invoke(args ...any) jsext.Value {
	return jsext.Value(e.Value().Invoke(args...))
}
func (e *Element) IsNaN() bool {
	return e.Value().IsNaN()
}
func (e *Element) IsNull() bool {
	return e.Value().IsNull()
}
func (e *Element) IsUndefined() bool {
	return e.Value().IsUndefined()
}
func (e *Element) Length() int {
	return e.Value().Length()
}
func (e *Element) Set(p string, x any) {
	e.Value().Set(p, x)
}
func (e *Element) SetIndex(i int, x any) {
	e.Value().SetIndex(i, x)
}
func (e *Element) String() string {
	return e.Value().String()
}
func (e *Element) Truthy() bool {
	return e.Value().Truthy()
}
func (e *Element) Type() js.Type {
	return e.Value().Type()
}
