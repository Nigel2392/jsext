package jse

import (
	"net/url"
	"syscall/js"

	"github.com/Nigel2392/jsext/v2"
)

type FormElement Element

// Form returns a form with the action, method, and id specified.
func Form(action, method, id string) *FormElement {
	var e = jsext.CreateElement("form")
	e.SetAttribute("action", action)
	e.SetAttribute("method", method)
	e.SetAttribute("id", id)
	return (*FormElement)(&e)
}

// JSValue returns the underlying js.Value.
func (e *FormElement) JSValue() js.Value {
	return js.Value(*e)
}

// Value returns as a Value(js.Value) wrapper.
func (e *FormElement) Element() *Element {
	return (*Element)(e)
}

// FormGroup returns a div with the classes specified.
func (e *FormElement) FormGroup(classes ...string) *FormElement {
	var l = Div(classes...)
	e.Element().AppendChild(l)
	return (*FormElement)(l)
}

// Label returns a label with the forElement and text specified.
func (e *FormElement) Label(forElement, text string, classes ...string) *FormElement {
	var l = Label(forElement, text, classes...)
	e.AppendChild(l)
	return (*FormElement)(l)
}

// Input returns an input with the type, name, placeholder, and value specified.
func (e *FormElement) Input(t, name string, opts *InputOptions) *FormElement {
	var i = Input(t, name, opts)
	e.AppendChild(i)
	return (*FormElement)(i)
}

// TimeInput returns an input to edit time with the name, placeholder, and value specified.
func (e *FormElement) TimeInput(name string, h, m, s int, opts *InputOptions) *Timer {
	var inp = TimeInput(name, h, m, s, opts)
	e.AppendChild(inp.Value())
	return inp
}

// Button returns a button with the text and classes specified.
func (e *FormElement) Button(innerText string, onClick func(this *Element, event jsext.Event)) *FormElement {
	var b = Button(innerText, onClick)
	e.AppendChild(b)
	return (*FormElement)(b)
}

// TextArea returns a textarea with the name, placeholder, and value specified.
func (e *FormElement) TextArea(name string, opts *InputOptions) *FormElement {
	var t = TextArea(name, opts)
	e.AppendChild(t)
	return (*FormElement)(t)
}

func (e *FormElement) AppendChild(children ...*Element) *FormElement {
	e.Element().AppendChild(children...)
	return e
}

func (e *FormElement) ID(s string) *FormElement {
	e.SetAttr("id", s)
	return e
}

func (e *FormElement) ClassList(classes ...string) jsext.Value {
	return e.Element().ClassList(classes...)
}

func (e *FormElement) SetAttr(key, value string) *FormElement {
	e.Element().SetAttr(key, value)
	return e
}

func (e *FormElement) DelAttr(key string) *FormElement {
	e.Element().DelAttr(key)
	return e
}

// OnSubmit sets the onsubmit event handler.
//
// This function will do nothing if the element on which this was called is not a html form.
func (e *FormElement) OnSubmit(autoclear bool, f func(this *Element, event jsext.Event, v url.Values)) *FormElement {

	var nodeName = e.Element().Get("nodeName")
	if nodeName.IsUndefined() || nodeName.IsNull() {
		return e
	}

	if nodeName.String() != "FORM" {
		return e
	}

	var newF = func(this *Element, event jsext.Event) {
		var formValues = make(map[string][]string)
		var form = event.Target()
		var elements = form.Get("elements")
		for i := 0; i < elements.Length(); i++ {
			var element = elements.Index(i)
			if element.IsUndefined() || element.IsNull() {
				continue
			}
			var name = element.Get("name").String()
			var value = element.Get("value").String()
			var mapValue, ok = formValues[name]
			if !ok {
				mapValue = make([]string, 0)
			}
			if autoclear {
				element.Set("value", "")
			}
			formValues[name] = append(mapValue, value)
		}
		f(this, event, formValues)
	}

	e.Element().AddEventListener("submit", newF)
	return e
}

// /////////////////////////////////////////////////////////
//
// js.Value methods.
//
// /////////////////////////////////////////////////////////
func (e *FormElement) Bool() bool {
	return e.JSValue().Bool()
}
func (e *FormElement) Call(m string, args ...any) jsext.Value {
	return jsext.Value(e.JSValue().Call(m, args...))
}
func (e *FormElement) Delete(p string) {
	e.JSValue().Delete(p)
}
func (e *FormElement) Equal(other js.Value) bool {
	return e.JSValue().Equal(other)
}
func (e *FormElement) Float() float64 {
	return e.JSValue().Float()
}
func (e *FormElement) Get(p string) jsext.Value {
	return jsext.Value(e.JSValue().Get(p))
}
func (e *FormElement) Index(i int) jsext.Value {
	return jsext.Value(e.JSValue().Index(i))
}
func (e *FormElement) Int() int {
	return e.JSValue().Int()
}
func (e *FormElement) Invoke(args ...any) jsext.Value {
	return jsext.Value(e.JSValue().Invoke(args...))
}
func (e *FormElement) IsNaN() bool {
	return e.JSValue().IsNaN()
}
func (e *FormElement) IsNull() bool {
	return e.JSValue().IsNull()
}
func (e *FormElement) IsUndefined() bool {
	return e.JSValue().IsUndefined()
}
func (e *FormElement) Length() int {
	return e.JSValue().Length()
}
func (e *FormElement) Set(p string, x any) {
	e.JSValue().Set(p, x)
}
func (e *FormElement) SetIndex(i int, x any) {
	e.JSValue().SetIndex(i, x)
}
func (e *FormElement) String() string {
	return e.JSValue().String()
}
func (e *FormElement) Truthy() bool {
	return e.JSValue().Truthy()
}
func (e *FormElement) Type() js.Type {
	return e.JSValue().Type()
}
