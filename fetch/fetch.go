package fetch

import (
	"io"
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/errs"
	"github.com/Nigel2392/jsext/v2/reader"
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

	ac := js.Global().Get("AbortController")
	if !ac.IsUndefined() {
		// Some browsers that support WASM don't necessarily support
		// the AbortController. See
		// https://developer.mozilla.org/en-US/docs/Web/API/AbortController#Browser_compatibility.
		ac = ac.New()
	}

	var jsReq, err = options.MarshalJS()
	if err != nil {
		return nil, err
	}

	var (
		fetchPromise     = js.Global().Call("fetch", options.URL, jsReq)
		respCh           = make(chan *Response, 1)
		errCh            = make(chan error, 1)
		success, failure js.Func
	)
	success = js.FuncOf(func(this js.Value, args []js.Value) any {
		success.Release()
		failure.Release()

		result := args[0]
		var jsHeaders = make(map[string][]string)
		var headersIt = result.Get("headers").Call("entries")
		for {
			n := headersIt.Call("next")
			if n.Get("done").Bool() {
				break
			}
			pair := n.Get("value")
			key, value := pair.Index(0).String(), pair.Index(1).String()
			jsHeaders[key] = append(jsHeaders[key], value)
		}

		var b = result.Get("body")
		var body io.ReadCloser
		// The body is undefined when the browser does not support streaming response bodies (Firefox),
		// and null in certain error cases, i.e. when the request is blocked because of CORS settings.
		if !b.IsUndefined() && !b.IsNull() {
			body = reader.NewStreamReader(b.Call("getReader"))
		} else {
			// Fall back to using ArrayBuffer
			// https://developer.mozilla.org/en-US/docs/Web/API/Body/arrayBuffer
			body = reader.NewArrayPromiseReader(result.Call("arrayBuffer"))
		}

		var code = result.Get("status").Int()
		respCh <- &Response{
			Status:     StatusText(code),
			StatusCode: code,
			Headers:    jsHeaders,
			Body:       body,
			Request:    &options,
			JS:         result,
		}

		return nil
	})
	failure = js.FuncOf(func(this js.Value, args []js.Value) any {
		success.Release()
		failure.Release()
		errCh <- errs.Error("jsext/fetch: Fetch() failed: " + args[0].Get("message").String())
		return nil
	})

	fetchPromise.Call("then", success, failure)
	select {
	case <-options.Context().Done():
		if !ac.IsUndefined() {
			// Abort the Fetch request.
			ac.Call("abort")
		}
		return nil, options.Context().Err()
	case resp := <-respCh:
		return resp, nil
	case err := <-errCh:
		return nil, err
	}
}
