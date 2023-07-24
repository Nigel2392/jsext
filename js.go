//go:build js && wasm
// +build js,wasm

package jsext

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"syscall/js"
	"time"
	"unsafe"

	"github.com/Nigel2392/jsext/v2/console"
	"github.com/Nigel2392/jsext/v2/errs"
	"github.com/Nigel2392/jsext/v2/export"
)

type Marshaller interface {
	MarshalJS() js.Value
}

type ErrorMarshaller interface {
	MarshalJS() (js.Value, error)
}

type Unmarshaller interface {
	UnmarshalJS(js.Value) error
}

type FuncMarshaller interface {
	MarshalJS() js.Func
}

type FuncUnmarshaller interface {
	UnmarshalJS(js.Func) error
}

// Default syscall/js values, some wrapped.
var (
	Export          export.Export
	Global          js.Value
	Document        js.Value
	DocumentValue   Value
	DocumentElement Element
	Window          Value
	Body            Element
	Head            Element
)

var waiter = make(chan struct{})

func Wait(exit chan error) error {
	select {
	case <-waiter:
		close(waiter)
		close(exit)
	case err := <-exit:
		close(waiter)
		close(exit)
		console.Log(err)
		return err
	}
	return nil
}

func StopWaiting() {
	waiter <- struct{}{}
}

func EmitInitiated() {
	EventEmit("jsext.initialized", strconv.FormatInt(time.Now().Unix(), 10))
}

func init() {
	// Initialize default values
	Global = js.Global()
	Document = Global.Get("document")
	DocumentValue = Value(Document)
	DocumentElement = Element(Document)
	Window = Value(Global.Get("window"))
	Body = Element(Document.Get("body"))
	Head = Element(Document.Get("head"))
	// Initialize jsext export object.
	Export = export.NewExport("jsext")
	Export.Set("runtime", Runtime.MarshalJS())
	// Register runtime eventlisteners.
	Runtime.Set("eventEmit", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		eventName := args[0].String()
		eventArgs := Args(args[1:])
		EventEmit(eventName, eventArgs.Slice()...)
		return nil
	}))
	Runtime.Set("eventOn", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var eventName = args[0].String()
		EventOn(eventName, func(a ...interface{}) {
			for _, arg := range args {
				if arg.Type() == js.TypeFunction {
					var Event = Global.Get("Event").New(eventName)
					Event.Set("args", a)
					arg.Invoke(Event)
				}
			}
		})
		return nil
	}))
}

type JSExtFunc func(this Value, args Args) interface{}

func (f JSExtFunc) MarshalJS() js.Func {
	var function = *(*func(this js.Value, args []js.Value) interface{})(unsafe.Pointer(&f))
	return js.FuncOf(function)
}

// Arguments for wrapped functions.
type Args []js.Value

// Len of arguments.
func (a Args) Len() int {
	return len(a)
}

// Event returns the first argument as an event.
func (a Args) Event() Event {
	if a.Len() == 0 {
		return Event(Undefined())
	}
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
			if v.InstanceOf(Global.Get("Array")) {
				s = append(s, ArrayToSlice(v))
			} else {
				s = append(s, v)
			}
		}
	}
	return s
}

// Register a function to the global window.
func RegisterFunc(name string, f FuncMarshaller) {
	Global.Set(name, f.MarshalJS())
}

// Get an element by id.
func GetElementById(id string) (Element, error) {
	var v = Document.Call("getElementById", id)
	if v.IsUndefined() || v.IsNull() {
		return Element(Undefined()), errs.Error("value is undefined || null")
	}
	return Element(v), nil
}

// Get an element by tag name.
func GetElementsByTagName(tag string) (Elements, error) {
	var v = Document.Call("getElementsByTagName", tag)
	if v.IsUndefined() || v.IsNull() {
		return nil, errs.Error("value is undefined || null")
	}
	var elements = UnpackArray[Element](v)
	return elements, nil
}

// Get an element by class name.
func GetElementsByClassName(class string) (Elements, error) {
	var v = Document.Call("getElementsByClassName", class)
	if v.IsUndefined() || v.IsNull() {
		return nil, errs.Error("value is undefined || null")
	}
	var elements = UnpackArray[Element](v)
	return elements, nil
}

// QuerySelector returns the first element that matches the specified selector.
func QuerySelector(selector string) (Element, error) {
	var v = Document.Call("querySelector", selector)
	if v.IsUndefined() || v.IsNull() {
		return Element(Undefined()), errs.Error("value is undefined || null")
	}
	return Element(v), nil
}

// QuerySelectorAll returns all elements that match the specified selector.
func QuerySelectorAll(selector string) (Elements, error) {
	var v = Document.Call("querySelectorAll", selector)
	if v.IsUndefined() || v.IsNull() {
		return nil, errs.Error("value is undefined || null")
	}
	var elements = UnpackArray[Element](v)
	return elements, nil
}

// Unpack a js.value array to a slice.
func UnpackArray[T Value | js.Value | Element | Style | Event](v js.Value) []T {
	var slice []T = make([]T, v.Length())
	for i := 0; i < v.Length(); i++ {
		slice[i] = T(v.Index(i))
	}
	return slice
}

// Set a favicon for the document.
func SetFavicon(url string) error {
	// Get the first link element in the head
	var link, err = QuerySelector("link[rel='icon']")
	if err != nil {
		link = CreateLink(map[string]string{
			"rel": "icon",
		})
	}
	var t string
	if url[len(url)-3:] == "ico" {
		t = "image/x-icon"
	} else {
		t = "image/png"
	}
	if link.Value().Truthy() {
		link.SetAttribute("href", url)
		link.SetAttribute("type", t)
		return nil
	}
	Head.AppendChild(CreateLink(map[string]string{
		"rel":  "icon",
		"type": t,
		"href": url,
	}))
	return nil
}

// Eval evaluates raw javascript code, returns the result as a js.Value.
func Eval(script string) Value {
	fmt.Println(Global)
	return Value(Global.Call("eval", script))
}

// Set a timeout on a function.
func SetTimeout(f FuncMarshaller, timeout int) Value {
	return Value(Global.Call("setTimeout", f.MarshalJS(), timeout))
}

// Set an interval on a function.
func SetInterval(f FuncMarshaller, timeout int) Value {
	return Value(Global.Call("setInterval", f.MarshalJS(), timeout))
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
	switch v := value.(type) {
	case Marshaller:
		return Value(v.MarshalJS())
	case ErrorMarshaller:
		var jsV, err = v.MarshalJS()
		if err != nil {
			return Value(js.Null())
		}
		return Value(jsV)
	case js.Value:
		return Value(v)
	case FuncMarshaller:
		return Value(v.MarshalJS().Value)
	default:
		return Value(js.ValueOf(v))
	}
}

// Returns a new object.
func NewObject() Value {
	return Value(Global.Get("Object").New())
}

// Returns a new array.
func NewArray() Value {
	return Value(Global.Get("Array").New())
}

// Returns a new date object.
func NewDate() Value {
	return Value(Global.Get("Date").New())
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
		return errs.Error("cookie length exceeds 4096 bytes")
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

// Check if the user agent is mobile.
func IsMobile() bool {
	var ua string
	if Window.Get("navigator").IsUndefined() || Window.Get("navigator").IsNull() {
		ua = Global.Get("navigator").Get("userAgent").String()
	} else {
		ua = Window.Get("navigator").Get("userAgent").String()
	}
	// var regex string = "(Android|webOS|iP(home|od|ad)|BlackBerry|IEMobile|Opera Mini)"
	var regex = `/Mobile|iP(hone|od)|Android|BlackBerry|IEMobile|Kindle|Silk-Accelerated|(hpw|web)OS|Opera M(obi|ini)/`
	var re = regexp.MustCompile(regex)
	return re.MatchString(ua)
}

// Convert a js.Value to a map[string]T.
//
// Uses unsafe.Pointer to convert the value to the type of T.
//
// Example:
//
//	var m = ObjectToMapT[string](js.ValueOf(map[string]interface{
//		"hello": "world",
//	}))
//	fmt.Println(m["hello"])
func ObjectToMapT[T any](obj js.Value) map[string]T {
	var m = make(map[string]T)
	var keys = obj.Call("keys")
	for i := 0; i < keys.Length(); i++ {
		var key = keys.Index(i).String()
		switch any(*new(T)).(type) {
		case string:
			var v = obj.Get(key).String()
			m[key] = *(*T)(unsafe.Pointer(&v))
		case int, int8, int16, int32, int64:
			var v T
			switch any(*new(T)).(type) {
			case int:
				var intie = obj.Get(key).Int()
				v = *(*T)(unsafe.Pointer(&intie))
			case int8:
				var intie = int8(obj.Get(key).Int())
				v = *(*T)(unsafe.Pointer(&intie))
			case int16:
				var intie = int16(obj.Get(key).Int())
				v = *(*T)(unsafe.Pointer(&intie))
			case int32:
				var intie = int32(obj.Get(key).Int())
				v = *(*T)(unsafe.Pointer(&intie))
			case int64:
				var intie = int64(obj.Get(key).Int())
				v = *(*T)(unsafe.Pointer(&intie))
			}
			m[key] = v
		case float64, float32:
			var v T
			switch any(*new(T)).(type) {
			case float64:
				var floatie = obj.Get(key).Float()
				v = *(*T)(unsafe.Pointer(&floatie))
			case float32:
				var floatie = float32(obj.Get(key).Float())
				v = *(*T)(unsafe.Pointer(&floatie))
			}
			m[key] = v
		case bool:
			var v = obj.Get(key).Bool()
			m[key] = *(*T)(unsafe.Pointer(&v))
		case js.Value, Value, Element, Event:
			m[key] = any(obj.Get(key)).(T)
		case []byte:
			var b = make([]byte, obj.Get(key).Length())
			js.CopyBytesToGo(b, obj.Get(key))
			m[key] = any(b).(T)
		case []any:
			m[key] = any(ArrayToSlice(obj.Get(key))).(T)
		case map[string]any:
			m[key] = any(ObjectToMap(obj.Get(key))).(T)
		default:
			panic("unsupported type")
		}
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
			if v.InstanceOf(Global.Get("Array")) {
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
			if v.InstanceOf(Global.Get("Array")) {
				s[i] = ArrayToSlice(v)
			} else {
				s[i] = v
			}
		}
	}
	return s
}

// Convert a slice to a js.Value array.
func SliceToArray(s []any) Value {
	var arr = NewArray()
	for i, v := range s {
		//lint:ignore S1034 Need the switch statement as is.
		switch v.(type) {
		case map[string]interface{}:
			arr.SetIndex(i, MapToObject(v.(map[string]interface{})).Value())
		case []interface{}:
			arr.SetIndex(i, SliceToArray(v.([]interface{})).Value())
		case []byte:
			arr.SetIndex(i, BytesToArray(v.([]byte)).Value())
		default:
			arr.SetIndex(i, ValueOf(v).Value())
		}
	}
	return arr
}

// Convert a map to a js.Value object.
func MapToObject(m map[string]interface{}) Value {
	var obj = NewObject()
	for k, v := range m {
		//lint:ignore S1034 Need the switch statement as is.
		switch v.(type) {
		case map[string]interface{}:
			obj.Set(k, MapToObject(v.(map[string]interface{})).Value())
		case []interface{}:
			obj.Set(k, SliceToArray(v.([]interface{})).Value())
		case []byte:
			obj.Set(k, BytesToArray(v.([]byte)).Value())
		default:
			obj.Set(k, ValueOf(v).Value())
		}
	}
	return obj
}

// Convert a byte slice to a js.Value array.
func BytesToArray(b []byte) Value {
	var buffer js.Value = Global.Get("ArrayBuffer").New(len(b))
	js.CopyBytesToJS(buffer, b)
	return Value(buffer)
}

func Alert(message string) {
	Call("alert", message)
}

// Get a value from the global scope.
func Get(key string) Value {
	return Value(Global.Get(key))
}

// Set a value in the global scope.
func Set(key string, value any) {
	value = MarshallableArguments(value)[0]
	Global.Set(key, value)
}

// Call a function in the global scope.
func Call(key string, args ...any) Value {
	args = MarshallableArguments(args...)
	return Value(Global.Call(key, args...))
}

// New a value in the global scope.
func New(key string, args ...any) Value {
	args = MarshallableArguments(args...)
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
