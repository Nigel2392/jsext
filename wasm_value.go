//go:build js && wasm
// +build js,wasm

package jsext

import "syscall/js"

// Wrapper for javascript values to make life easier.
type Value js.Value

// Returns the js.Value of the value.
func (w Value) Value() js.Value {
	return js.Value(w)
}

// Returns the jsext.Element of the value.
func (w Value) ToElement() Element {
	return Element(w.Value())
}

func (w Value) Invalid() bool {
	return w.IsUndefined() || w.IsNull()
}

func (w Value) IsArray() bool {
	var arr = Global.Get("Array")
	var isArr = arr.Get("isArray")
	return isArr.Invoke(w.Value()).Bool()
}

func (w Value) IsObject() bool {
	return w.Value().Type() == js.TypeObject
}

func (w Value) IsString() bool {
	return w.Value().Type() == js.TypeString
}

func (w Value) IsNumber() bool {
	return w.Value().Type() == js.TypeNumber
}

func (w Value) IsBoolean() bool {
	return w.Value().Type() == js.TypeBoolean
}

func (w Value) IsFunction() bool {
	return w.Value().Type() == js.TypeFunction
}

// Returns true if the value is a javascript Element or HTMLElement.
func (w Value) IsElement() bool {
	var HTMLElement = js.Global().Get("HTMLElement")
	var Element = js.Global().Get("Element")
	return InstanceOf(w.Value(), HTMLElement) || InstanceOf(w.Value(), Element)
}

// Convert a slice of strings to a slice of interfaces.
func SliceToInterface(s []string) []interface{} {
	var anyList = make([]interface{}, len(s))
	for i, v := range s {
		anyList[i] = v
	}
	return anyList
}

// Toggle an attribute on an element.
func (w Value) Toggle(s string) {
	w.Value().Call("toggle", s)
}

// Great for values which cannot be addressed with .toggle()
// 	- Returns the new value
func (w Value) ToggleCustom(s string) bool {
	var val = w.Value().Get(s)
	var nVal = !val.Bool()
	w.Value().Set(s, nVal)
	return nVal
}

// Add attributes to an element.
func (w Value) Add(s ...string) {
	w.Value().Call("add", SliceToInterface(s)...)
}

// Remove an attribute from an element, or the entire element if no attribute is given.
func (w Value) Remove(a ...string) {
	if len(a) == 0 {
		w.Value().Call("remove")
	} else {
		w.Value().Call("remove", SliceToInterface(a)...)
	}
}

// AppendChild inserts a node after another node.
func (w Value) AppendChild(child Element) {
	// if !child.IsElement() || !w.IsElement() {
	// panic("replaceChild: child and w must be elements")
	// }
	w.Value().Call("appendChild", child.Value())
}

// PrependChild inserts a node before another node.
func (w Value) PrependChild(child Element) {
	// if !child.IsElement() || !w.IsElement() {
	// panic("replaceChild: child and w must be elements")
	// }
	w.Value().Call("prepend", child.Value())
}

// InsertBefore inserts a node before another node.
func (w Value) InsertBefore(child Element, before Element) {
	// if !before.IsElement() || !child.IsElement() || !w.IsElement() {
	// panic("replaceChild: before, child and w must be elements")
	// }
	w.Value().Call("insertBefore", child.Value(), before.Value())
}

// ReplaceChild replaces the child with the before element.
func (w Value) ReplaceChild(child Element, before Element) {
	// if !before.IsElement() || !child.IsElement() || !w.IsElement() {
	// panic("replaceChild: before, child and w must be elements")
	// }
	w.Value().Call("replaceChild", child.Value(), before.Value())
}

// GetElementById returns the first element with the given id.
func (w Value) GetElementById(id string) Element {
	return Element(w.Call("getElementById", id))
}

// GetElementByTagName returns the first element with the given tag name.
func (w Value) GetElementByTagName(tag string) Element {
	return Element(w.Call("getElementsByTagName", tag))
}

// GetElementByClassName returns the first element with the given class name.
func (w Value) GetElementByClassName(class string) Element {
	return Element(w.Call("getElementsByClassName", class))
}

// Query select inside of the value.
func (w Value) QuerySelector(selector string) Element {
	return Element(w.Call("querySelector", selector))
}

// Query select all inside of the value.
func (w Value) QuerySelectorAll(selector string) Elements {
	var els = w.Call("querySelectorAll", selector)
	var elements []Element = make([]Element, els.Length())
	for i := 0; i < els.Length(); i++ {
		elements[i] = Element(els.Index(i))
	}
	return elements
}

///////////////////////////////////////////////////////////
//
// js.Value methods.
//
///////////////////////////////////////////////////////////
func (w Value) Bool() bool {
	return w.Value().Bool()
}
func (w Value) Call(m string, args ...any) Value {
	return Value(w.Value().Call(m, args...))
}
func (w Value) Delete(p string) {
	w.Value().Delete(p)
}
func (w Value) Equal(other js.Value) bool {
	return w.Value().Equal(other)
}
func (w Value) Float() float64 {
	return w.Value().Float()
}
func (w Value) Get(p string) Value {
	return Value(w.Value().Get(p))
}
func (w Value) Index(i int) Value {
	return Value(w.Value().Index(i))
}
func (w Value) InstanceOf(t js.Value) bool {
	return w.Value().InstanceOf(t)
}
func (w Value) Int() int {
	return w.Value().Int()
}
func (w Value) Invoke(args ...any) Value {
	return Value(w.Value().Invoke(args...))
}
func (w Value) IsNaN() bool {
	return w.Value().IsNaN()
}
func (w Value) IsNull() bool {
	return w.Value().IsNull()
}
func (w Value) IsUndefined() bool {
	return w.Value().IsUndefined()
}
func (w Value) Length() int {
	return w.Value().Length()
}
func (w Value) New(args ...any) Value {
	return Value(w.Value().New(args...))
}
func (w Value) Set(p string, x any) {
	w.Value().Set(p, x)
}
func (w Value) SetIndex(i int, x any) {
	w.Value().SetIndex(i, x)
}
func (w Value) String() string {
	return w.Value().String()
}
func (w Value) Truthy() bool {
	return w.Value().Truthy()
}
func (w Value) Type() js.Type {
	return w.Value().Type()
}
