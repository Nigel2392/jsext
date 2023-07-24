package xhr

import "syscall/js"

type ProgressEvent js.Value

func (p ProgressEvent) MarshalJS() js.Value {
	return js.Value(p)
}

func (p ProgressEvent) LengthComputable() bool {
	return js.Value(p).Get("lengthComputable").Bool()
}

func (p ProgressEvent) Loaded() int {
	return js.Value(p).Get("loaded").Int()
}

func (p ProgressEvent) Total() int {
	return js.Value(p).Get("total").Int()
}

func (p ProgressEvent) Percent() float64 {
	return float64(p.Loaded()) / float64(p.Total())
}

func (p ProgressEvent) Target() js.Value {
	return js.Value(p).Get("target")
}

func (p ProgressEvent) Set(key string, value interface{}) {
	js.Value(p).Set(key, value)
}

func (p ProgressEvent) Get(key string) js.Value {
	return js.Value(p).Get(key)
}

func (p ProgressEvent) Call(key string, args ...interface{}) js.Value {
	return js.Value(p).Call(key, args...)
}
