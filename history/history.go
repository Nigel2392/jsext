package history

import (
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/jsc"
)

type History js.Value

func Get() History {
	return History(js.Global().Get("history"))
}

func (h History) JSValue() js.Value {
	return js.Value(h)
}

func (h History) Get(key string) js.Value {
	return js.Value(h).Get(key)
}

func (h History) Call(method string, args ...interface{}) (js.Value, error) {
	for i, arg := range args {
		var v, err = jsc.ValueOf(arg)
		if err != nil {
			return js.Value{}, err
		}
		args[i] = v
	}
	return js.Value(h).Call(method, args...), nil
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

func (h History) PushState(data interface{}, title string, url string) error {
	var v, err = jsc.ValueOf(data)
	if err != nil {
		return err
	}
	h.Call("pushState", v, title, url)
	return nil
}

func (h History) ReplaceState(data interface{}, title string, url string) error {
	var v, err = jsc.ValueOf(data)
	if err != nil {
		return err
	}
	h.Call("replaceState", v, title, url)
	return nil
}

func (h History) State() js.Value {
	return h.Get("state")
}

func (h History) ScanState(dst interface{}) error {
	return jsc.Scan(h.State(), dst)
}
