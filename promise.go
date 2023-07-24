package jsext

import "syscall/js"

type Promise struct {
	// The underlying javascript value of the Promise.
	p Value
}

// Returns the jsext.Value of the Promise.
func (w Promise) Value() Value {
	return w.p
}

// Returns the js.Value of the Promise.
func (w Promise) MarshalJS() js.Value {
	return w.p.Value()
}

// Create a new Promise that resolves to the given value.
func NewPromiseResolve(value any) Promise {
	var promiseConstructor = js.Global().Get("Promise")
	var promise = promiseConstructor.New(js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return value
	}))
	return Promise{
		p: Value(promise),
	}
}

// Create a new Promise.
func NewPromise(f func() (any, error)) Promise {
	var fn = func(this js.Value, args []js.Value) interface{} {
		var resolve = args[0]
		var reject = args[1]
		go func() {
			var result, err = f()
			if err != nil {
				var jsErrConstructor = js.Global().Get("Error")
				var jsErr = jsErrConstructor.New(err.Error())
				reject.Invoke(jsErr)
			} else {
				resolve.Invoke(ValueOf(result).Value())
			}
		}()
		return nil
	}
	var promiseConstructor = js.Global().Get("Promise")
	var promise = promiseConstructor.New(js.FuncOf(fn))
	return Promise{Value(promise)}
}

// Then a function on the Promise.
func (w Promise) Then(f func(Value)) Promise {
	var fn = func(this js.Value, args []js.Value) interface{} {
		if len(args) == 0 {
			return nil
		}
		var result = args[0]
		f(Value(result))
		return nil
	}
	var promise = w.MarshalJS().Call("then", js.FuncOf(fn))
	return Promise{Value(promise)}
}

// Catch an error from the Promise.
func (w Promise) Catch(f func(error)) Promise {
	var fn = func(this js.Value, args []js.Value) interface{} {
		if len(args) == 0 {
			return nil
		}
		var err = args[0]
		f(js.Error{Value: err})
		return nil
	}
	var promise = w.MarshalJS().Call("catch", js.FuncOf(fn))
	return Promise{Value(promise)}
}
