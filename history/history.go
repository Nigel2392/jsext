package history

import (
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/jsc"
)

var history = js.Global().Get("history")

func Call(method string, args ...interface{}) (js.Value, error) {
	for i, arg := range args {
		var v, err = jsc.ValueOf(arg)
		if err != nil {
			return js.Value{}, err
		}
		args[i] = v
	}
	return js.Value(history).Call(method, args...), nil
}

func Length() int {
	return history.Get("length").Int()
}

func Go(index int) {
	history.Call("go", index)
}

func Back() {
	history.Call("back")
}

func Forward() {
	history.Call("forward")
}

func PushState(data interface{}, title string, url string) error {
	var v, err = jsc.ValueOf(data)
	if err != nil {
		return err
	}
	history.Call("pushState", v, title, url)
	return nil
}

func ReplaceState(data interface{}, title string, url string) error {
	var v, err = jsc.ValueOf(data)
	if err != nil {
		return err
	}
	history.Call("replaceState", v, title, url)
	return nil
}

func State() js.Value {
	return history.Get("state")
}

func ScanState(dst interface{}) error {
	return jsc.Scan(State(), dst)
}
