package history

import (
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/jsc"
)

type History js.Value

func GetHistory() History {
	return History(js.Global().Get("history"))
}

func (h History) JSValue() js.Value {
	return js.Value(h)
}

func (h History) Get(key string) js.Value {
	return js.Value(h).Get(key)
}

func (h History) Call(method string, args ...interface{}) js.Value {
	for i, arg := range args {
		args[i] = jsc.ValueOf(arg)
	}
	return js.Value(h).Call(method, args...)
}

func (h History) Length() int {
	return h.Get("length").Int()
}

func (h History) Go(index int) {
	h.Call("go", index)
}

func (h History) Back() {
	h.Call("back")
}

func (h History) Forward() {
	h.Call("forward")
}

func (h History) PushState(data interface{}, title string, url string) {
	h.Call("pushState", jsc.ValueOf(data), title, url)
}

func (h History) ReplaceState(data interface{}, title string, url string) {
	h.Call("replaceState", jsc.ValueOf(data), title, url)
}

func (h History) State() js.Value {
	return h.Get("state")
}

func (h History) ScanState(dst interface{}) error {
	return jsc.Scan(h.State(), dst)
}
