package history

import (
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/jsc"
)

var history = js.Global().Get("history")

// Call a method on the history object.
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

// Get the length of the history object.
func Length() int {
	return history.Get("length").Int()
}

// Go to a specific page in the history.
func Go(index int) {
	history.Call("go", index)
}

// Go back one page.
func Back() {
	history.Call("back")
}

// Go forward one page.
func Forward() {
	history.Call("forward")
}

// Push a state to the history.
func PushState(data interface{}, title string, url string) error {
	var v, err = jsc.ValueOf(data)
	if err != nil {
		return err
	}
	history.Call("pushState", v, title, url)
	return nil
}

// Replace the current state in the history.
func ReplaceState(data interface{}, title string, url string) error {
	var v, err = jsc.ValueOf(data)
	if err != nil {
		return err
	}
	history.Call("replaceState", v, title, url)
	return nil
}

// Get the current state.
func State() js.Value {
	return history.Get("state")
}

// Scan the current state into dst.
func ScanState(dst interface{}) error {
	return jsc.Scan(State(), dst)
}
