//go:build js && wasm
// +build js,wasm

package jsext

import (
	"syscall/js"
)

var (
	JSExt    Export
	Global   js.Value
	Document js.Value
	Window   js.Value
	Body     Element
	Head     Element
)

func init() {
	Global = js.Global()
	Document = Global.Get("document")
	Window = Global.Get("window")
	Body = Element(Document.Get("body"))
	Head = Element(Document.Get("head"))
	JSExt = NewExport()
	JSExt.Register("jsext")
}

type JSFunc func(this js.Value, args []js.Value) interface{}
type JSExtFunc func(this Value, args Args) interface{}
type JSExtEventFunc func(this Value, event Event) interface{}

func ToJSFunc(f JSExtFunc) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return f(Value(this), Args(args))
	})
}

func (f JSExtFunc) ToJSFunc() js.Func {
	return ToJSFunc(f)
}

func (f JSExtFunc) FromJSFunc(fun JSFunc) JSExtFunc {
	return JSExtFunc(func(this Value, args Args) interface{} {
		return fun(this.Value(), args.Value())
	})
}

type Args []js.Value

func (a Args) Len() int {
	return len(a)
}

func (a Args) Event() Event {
	return Event(a[0])
}

func (a Args) Value() []js.Value {
	return []js.Value(a)
}

func RegisterFunc(name string, f JSExtFunc) {
	js.Global().Set(name, f.ToJSFunc())
}

func WrapFunc(f func()) js.Func {
	return JSExtFunc(func(this Value, args Args) interface{} {
		f()
		return nil
	}).ToJSFunc()
}

func GetElementById(id string) Element {
	return Element(Document.Call("getElementById", id))
}

func GetElementByTagName(tag string) Element {
	return Element(Document.Call("getElementsByTagName", tag))
}

func GetElementByClassName(class string) Element {
	return Element(Document.Call("getElementsByClassName", class))
}

func QuerySelector(selector string) Element {
	return Element(Document.Call("querySelector", selector))
}

func QuerySelectorAll(selector string) Elements {
	var els = Document.Call("querySelectorAll", selector)
	var elements []Element = make([]Element, els.Length())
	for i := 0; i < els.Length(); i++ {
		elements[i] = Element(els.Index(i))
	}
	return elements
}

func SetFavicon(url string) {
	// Get the first link element in the head
	var link = QuerySelector("link[rel='icon']")
	var t string
	if url[len(url)-3:] == "ico" {
		t = "image/x-icon"
	} else {
		t = "image/png"
	}
	if link.Value().Truthy() {
		link.SetAttribute("href", url)
		link.SetAttribute("type", t)
		return
	}
	Head.AppendChild(CreateLink(map[string]string{
		"rel":  "icon",
		"type": t,
		"href": url,
	}))
}

func Eval(script string) Value {
	return Value(js.Global().Call("eval", script))
}

func AddEventListenerById(id, event string, listener func(this Value, event Event)) {
	GetElementById(id).AddEventListener(event, listener)
}

func RemoveEventListenerById(id, event string, listener func(this Value, event Event)) {
	GetElementById(id).RemoveEventListener(event, listener)
}

func SetTimeout(f JSExtFunc, timeout int) Value {
	return Value(js.Global().Call("setTimeout", f.ToJSFunc(), timeout))
}

func SetInterval(f JSExtFunc, timeout int) Value {
	return Value(js.Global().Call("setInterval", f.ToJSFunc(), timeout))
}

func CreateElement(tag string) Element {
	return Element(Document.Call("createElement", tag))
}

func CreateLink(kv map[string]string) Element {
	link := CreateElement("link")
	for k, v := range kv {
		link.SetAttribute(k, v)
	}
	return Element(link)
}

func InstanceOf(value, constructor js.Value) bool {
	return value.InstanceOf(constructor)
}

func TypeOf(value, constructor js.Value) bool {
	return value.Type() == constructor.Type()
}

func ValueOf(value any) Value {
	return Value(js.ValueOf(value))
}

func NewObject() Value {
	return Value(js.Global().Get("Object").New())
}

func NewArray() Value {
	return Value(js.Global().Get("Array").New())
}

func NewDate() Value {
	return Value(js.Global().Get("Date").New())
}

//type JavaScript interface {
//	JSValue() js.Value
//	Value() Value
//	Bool() bool
//	Call(m string, args ...any) js.Value
//	Delete(p string)
//	Equal(w js.Value) bool
//	Float() float64
//	Get(p string) js.Value
//	Index(i int) js.Value
//	InstanceOf(t js.Value) bool
//	Int() int
//	Invoke(args ...any) js.Value
//	IsNaN() bool
//	IsNull() bool
//	IsUndefined() bool
//	Length() int
//	New(args ...any) js.Value
//	Set(p string, x any)
//	SetIndex(i int, x any)
//	String() string
//	Truthy() bool
//	Type() js.Type
//}
