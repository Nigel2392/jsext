//go:build js && wasm
// +build js,wasm

package jsext

import (
	"errors"
	"strings"
	"syscall/js"
	"time"
)

// Default syscall/js values, some wrapped.
var (
	JSExt    Export
	Global   js.Value
	Document js.Value
	Window   js.Value
	Body     Element
	Head     Element
)

func init() {
	// Initialize default values
	Global = js.Global()
	Document = Global.Get("document")
	Window = Global.Get("window")
	Body = Element(Document.Get("body"))
	Head = Element(Document.Get("head"))
	// Initialize jsext export object.
	JSExt = NewExport()
	JSExt.Register("jsext")
	// Register runtime eventlisteners.
	Runtime.RegisterTo("runtime", JSExt.JSExt())
	Runtime.SetFuncWithArgs("eventEmit", func(this Value, args Args) interface{} {
		eventName := args[0].String()
		eventArgs := args[1:]
		EventEmit(eventName, eventArgs.Slice()...)
		return nil
	})
	Runtime.SetFuncWithArgs("eventOn", func(this Value, args Args) interface{} {
		eventName := args[0].String()
		EventOn(eventName, func(a ...interface{}) {
			for _, arg := range args {
				if arg.Type() == js.TypeFunction {
					var Event = js.Global().Get("Event").New(eventName)
					Event.Set("args", a)
					arg.Invoke(Event)
				}
			}
		})
		return nil
	})

}

// Default functions, wrapped.
type JSFunc func(this js.Value, args []js.Value) interface{}
type JSExtFunc func(this Value, args Args) interface{}
type JSExtEventFunc func(this Value, event Event) interface{}

// Conver to javascript function.
func ToJSFunc(f JSExtFunc) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return f(Value(this), Args(args))
	})
}

// Convert to javascript function.
func (f JSExtFunc) ToJSFunc() js.Func {
	return ToJSFunc(f)
}

// Convert from javascript function.
func (f JSExtFunc) FromJSFunc(fun JSFunc) JSExtFunc {
	return JSExtFunc(func(this Value, args Args) interface{} {
		return fun(this.Value(), args.Value())
	})
}

// Arguments for wrapped functions.
type Args []js.Value

// Len of arguments.
func (a Args) Len() int {
	return len(a)
}

// Event returns the first argument as an event.
func (a Args) Event() Event {
	return Event(a[0])
}

// Value returns the arguments as a slice of js.Value.
func (a Args) Value() []js.Value {
	return []js.Value(a)
}

// Return a slice of any.
func (a Args) Slice() []interface{} {
	var s = make([]interface{}, 0)
	for _, v := range a {
		switch v.Type() {
		case js.TypeObject:
			s = append(s, ObjectToMap(v))
		case js.TypeNull:
			s = append(s, nil)
		case js.TypeBoolean:
			s = append(s, v.Bool())
		case js.TypeNumber:
			s = append(s, v.Float())
		case js.TypeString:
			s = append(s, v.String())
		default:
			if v.InstanceOf(js.Global().Get("Array")) {
				s = append(s, ArrayToSlice(v))
			} else {
				s = append(s, v)
			}
		}
	}
	return s
}

// Register a function to the global window.
func RegisterFunc(name string, f JSExtFunc) {
	js.Global().Set(name, f.ToJSFunc())
}

// Wrap a function, convert it to a js.Func.
func WrapFunc(f func()) js.Func {
	return JSExtFunc(func(this Value, args Args) interface{} {
		f()
		return nil
	}).ToJSFunc()
}

// Get an element by id.
func GetElementById(id string) Element {
	return Element(Document.Call("getElementById", id))
}

// Get an element by tag name.
func GetElementByTagName(tag string) Element {
	return Element(Document.Call("getElementsByTagName", tag))
}

// Get an element by class name.
func GetElementByClassName(class string) Element {
	return Element(Document.Call("getElementsByClassName", class))
}

// QuerySelector returns the first element that matches the specified selector.
func QuerySelector(selector string) Element {
	return Element(Document.Call("querySelector", selector))
}

// QuerySelectorAll returns all elements that match the specified selector.
func QuerySelectorAll(selector string) Elements {
	var els = Document.Call("querySelectorAll", selector)
	var elements []Element = make([]Element, els.Length())
	for i := 0; i < els.Length(); i++ {
		elements[i] = Element(els.Index(i))
	}
	return elements
}

// Set a favicon for the document.
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

// Eval evaluates raw javascript code, returns the result as a js.Value.
func Eval(script string) Value {
	return Value(js.Global().Call("eval", script))
}

// Add an event listener to an element by it's id.
func AddEventListenerById(id, event string, listener func(this Value, event Event)) {
	GetElementById(id).AddEventListener(event, listener)
}

// Remove an event listener from an element by it's id.
func RemoveEventListenerById(id, event string, listener func(this Value, event Event)) {
	GetElementById(id).RemoveEventListener(event, listener)
}

// Set a timeout on a function.
func SetTimeout(f JSExtFunc, timeout int) Value {
	return Value(js.Global().Call("setTimeout", f.ToJSFunc(), timeout))
}

// Set an interval on a function.
func SetInterval(f JSExtFunc, timeout int) Value {
	return Value(js.Global().Call("setInterval", f.ToJSFunc(), timeout))
}

// Create a javascript HTMLElement.
func CreateElement(tag string) Element {
	return Element(Document.Call("createElement", tag))
}

// Create a link element.
func CreateLink(kv map[string]string) Element {
	link := CreateElement("link")
	for k, v := range kv {
		link.SetAttribute(k, v)
	}
	return Element(link)
}

// Returns if something is an instance of something else.
func InstanceOf(value, constructor js.Value) bool {
	return value.InstanceOf(constructor)
}

// Returns if something is a type of something else.
func TypeOf(value, constructor js.Value) bool {
	return value.Type() == constructor.Type()
}

// Returns the value of a property.
func ValueOf(value any) Value {
	return Value(js.ValueOf(value))
}

// Returns a new object.
func NewObject() Value {
	return Value(js.Global().Get("Object").New())
}

// Returns a new array.
func NewArray() Value {
	return Value(js.Global().Get("Array").New())
}

// Returns a new date object.
func NewDate() Value {
	return Value(js.Global().Get("Date").New())
}

// Returns a new undefined value.
func Undefined() Value {
	return Value(js.Undefined())
}

// Returns a new null value.
func Null() Value {
	return Value(js.Null())
}

// Set a document cookie
func SetCookie(name, value string, tim time.Duration) error {
	var expires = time.Now().Add(tim).UTC().Format(time.RFC1123)
	var cookie = name + "=" + value + "; expires=" + expires + "; path=/"
	if len(cookie) > 4096 {
		return errors.New("cookie length exceeds 4096 bytes")
	}
	Document.Set("cookie", cookie)
	return nil
}

// Get a document cookie
func GetCookie(name string) string {
	var cookie = Document.Get("cookie").String()
	var parts = strings.Split(cookie, ";")
	for _, part := range parts {
		var kv = strings.Split(part, "=")
		if strings.TrimSpace(kv[0]) == name {
			return strings.TrimSpace(kv[1])
		}
	}
	return ""
}

// Delete a document cookie
func DeleteCookie(name string) {
	SetCookie(name, "", -1)
}

// Convert a js.Value to a map[string]string.
func ObjectToMapString(obj js.Value) map[string]string {
	var m = make(map[string]string)
	var keys = obj.Call("keys")
	for i := 0; i < keys.Length(); i++ {
		var key = keys.Index(i).String()
		m[key] = obj.Get(key).String()
	}
	return m
}

// Convert a js.Value to a map[string]interface{}.
func ObjectToMap(obj js.Value) map[string]interface{} {
	var m = make(map[string]interface{})
	var keys = obj.Call("keys")
	for i := 0; i < keys.Length(); i++ {
		var key = keys.Index(i).String()
		var v = obj.Get(key)
		switch v.Type() {
		case js.TypeObject:
			m[key] = ObjectToMap(v)
		case js.TypeNull:
			m[key] = nil
		case js.TypeBoolean:
			m[key] = v.Bool()
		case js.TypeNumber:
			m[key] = v.Float()
		case js.TypeString:
			m[key] = v.String()
		default:
			if v.InstanceOf(js.Global().Get("Array")) {
				m[key] = ArrayToSlice(v)
			} else {
				m[key] = v
			}
		}
	}
	return m
}

// Convert a js.Value array to a []interface{}.
func ArrayToSlice(arr js.Value) []interface{} {
	var s = make([]interface{}, arr.Length())
	for i := 0; i < arr.Length(); i++ {
		var v = arr.Index(i)
		switch v.Type() {
		case js.TypeObject:
			s[i] = ObjectToMap(v)
		case js.TypeNull:
			s[i] = nil
		case js.TypeBoolean:
			s[i] = v.Bool()
		case js.TypeNumber:
			s[i] = v.Float()
		case js.TypeString:
			s[i] = v.String()
		default:
			if v.InstanceOf(js.Global().Get("Array")) {
				s[i] = ArrayToSlice(v)
			} else {
				s[i] = v
			}
		}
	}
	return s
}

// Get a value from the global scope.
func Get(key string) Value {
	return Value(Global.Get(key))
}

// Set a value in the global scope.
func Set(key string, value any) {
	Global.Set(key, value)
}

// Call a function in the global scope.
func Call(key string, args ...any) Value {
	return Value(Global.Call(key, args...))
}

// New a value in the global scope.
func New(key string, args ...any) Value {
	return Value(Global.Get(key).New(args...))
}

// Delete a value in the global scope.
func Delete(key string) {
	Global.Delete(key)
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
