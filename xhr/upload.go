package xhr

import "syscall/js"

type XMLHttpRequestUpload struct {
	Value *js.Value
}

func (x *XMLHttpRequestUpload) MarshalJS() js.Value {
	return *x.Value
}

func (x *XMLHttpRequestUpload) Set(key string, value interface{}) {
	x.Value.Set(key, value)
}

func (x *XMLHttpRequestUpload) Abort(f func(event ProgressEvent)) js.Func {
	if f == nil {
		x.Set("onabort", nil)
		return js.Func{}
	}
	var fn = jsFuncFromEventFunc(f)
	x.Set("onabort", fn)
	return fn
}

func (x *XMLHttpRequestUpload) Error(f func(event ProgressEvent)) js.Func {
	if f == nil {
		x.Set("onerror", nil)
		return js.Func{}
	}
	var fn = jsFuncFromEventFunc(f)
	x.Set("onerror", fn)
	return fn
}

func (x *XMLHttpRequestUpload) Load(f func(event ProgressEvent)) js.Func {
	if f == nil {
		x.Set("onload", nil)
		return js.Func{}
	}
	var fn = jsFuncFromEventFunc(f)
	x.Set("onload", fn)
	return fn
}

func (x *XMLHttpRequestUpload) LoadStart(f func(event ProgressEvent)) js.Func {
	if f == nil {
		x.Set("onloadstart", nil)
		return js.Func{}
	}
	var fn = jsFuncFromEventFunc(f)
	x.Set("onloadstart", fn)
	return fn
}

func (x *XMLHttpRequestUpload) LoadEnd(f func(event ProgressEvent)) js.Func {
	if f == nil {
		x.Set("onloadend", nil)
		return js.Func{}
	}
	var fn = jsFuncFromEventFunc(f)
	x.Set("onloadend", fn)
	return fn
}

func (x *XMLHttpRequestUpload) Progress(f func(event ProgressEvent)) js.Func {
	if f == nil {
		x.Set("onprogress", nil)
		return js.Func{}
	}
	var fn = jsFuncFromEventFunc(f)
	x.Set("onprogress", fn)
	return fn
}

func (x *XMLHttpRequestUpload) Timeout(f func(event ProgressEvent)) js.Func {
	if f == nil {
		x.Set("ontimeout", nil)
		return js.Func{}
	}
	var fn = jsFuncFromEventFunc(f)
	x.Set("ontimeout", fn)
	return fn
}

func jsFuncFromEventFunc(f func(event ProgressEvent)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var event = ProgressEvent(args[0])
		f(event)
		return nil
	})
}
