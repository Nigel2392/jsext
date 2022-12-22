package jsext

import "syscall/js"

type Export Value

func NewExport() Export {
	return Export(js.Global().Get("Object").New())
}

func (e Export) Value() js.Value {
	return js.Value(e)
}

func (e Export) Set(name string, value interface{}) {
	js.Value(e).Set(name, value)
}

func (e Export) Get(name string) Value {
	return Value(js.Value(e).Get(name))
}

func (e Export) Call(name string, args ...interface{}) Value {
	return Value(js.Value(e).Call(name, args...))
}

func (e Export) SetFunc(name string, f func()) {
	js.Value(e).Set(name, WrapFunc(f))
}

func (e Export) SetFuncWithArgs(name string, f JSExtFunc) {
	js.Value(e).Set(name, f.ToJSFunc())
}

func (e Export) SetMultipleWithArgs(fns map[string]JSExtFunc) {
	for name, f := range fns {
		e.SetFuncWithArgs(name, f)
	}
}

func (e Export) Remove(name string) {
	js.Value(e).Delete(name)
}

func (e Export) Register(name string) {
	Global.Set(name, e.Value())
}

func (e Export) RegisterTo(name string, to Value) {
	to.Set(name, e.Value())
}

func (e Export) RegisterToExport(name string, to Export) {
	to.Set(name, e.Value())
}
