package fetch

import (
	"errors"
	"syscall/js"
)

// TinyGO fetch request implementation.
func Fetch(options *Request) (*Response, error) {
	return fetch(*options)
}

func fetch(options Request) (*Response, error) {
	if options.Method == "" {
		options.Method = "GET"
	}

	switch options.Method {
	case "GET", "HEAD":
		options.Body = nil
	case "":
		options.Method = "GET"
		options.Body = nil
	}

	var fetch = js.Global().Call("fetch", options.URL, options.MarshalJS())
	if fetch.IsUndefined() {
		panic("fetch is undefined")
	}

	var respChan = make(chan *Response)
	var then = fetch.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var response = args[0]
		var headers = response.Get("headers")
		var jsHeaders = make(map[string][]string)
		keys := headers.Call("keys")
		for i := 0; i < keys.Length(); i++ {
			var key = keys.Index(i).String()
			var value = headers.Call("get", key)
			var values []string
			if value.Type() == js.TypeString {
				values = []string{value.String()}
			} else {
				for j := 0; j < value.Length(); j++ {
					values = append(values, value.Index(j).String())
				}
			}
			jsHeaders[key] = values
		}
		var statusCode = response.Get("status").Int()
		var resp = &Response{
			Headers:    jsHeaders,
			StatusCode: statusCode,
			JS:         response,
		}
		response.Call("text").Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			var b []byte
			switch {
			case len(args) < 1:
				b = []byte{}
			case args[0].IsUndefined(), args[0].IsNull():
				b = []byte{}
			case args[0].Type() == js.TypeString:
				b = []byte(args[0].String())
			case js.Global().Get("Uint8Array").Call("isPrototypeOf", args[0]).Bool(),
				js.Global().Get("Uint8ClampedArray").Call("isPrototypeOf", args[0]).Bool(),
				args[0].InstanceOf(js.Global().Get("ArrayBuffer")):
				// Is a buffer type.
				var uint8Array = args[0]
				var length = uint8Array.Get("length").Int()
				b = make([]byte, length)
				js.CopyBytesToGo(b, uint8Array)
			default:
				b = []byte{}
			}
			resp.Body = []byte(b)
			respChan <- resp
			return nil
		}))
		return nil
	}))
	if then.IsUndefined() {
		close(respChan)
		return nil, errors.New("then is undefined")
	}
	var errChan = make(chan error)
	then.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var err = args[0]
		var errString = err.Get("message").String()
		errChan <- errors.New(errString)
		return nil
	}))
	var resp *Response
	var err error
	select {
	case resp = <-respChan:
	case err = <-errChan:
	}
	close(respChan)
	close(errChan)
	return resp, err
}
