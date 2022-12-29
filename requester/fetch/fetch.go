package fetch

import (
	"errors"
	"strings"
	"syscall/js"
)

type Request struct {
	Body        []byte            `json:"body,omitempty"`
	BodyUsed    bool              `json:"bodyUsed,omitempty"`
	Cache       string            `json:"cache,omitempty"`
	Credentials string            `json:"credentials,omitempty"`
	Destination string            `json:"destination,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Integrity   string            `json:"integrity,omitempty"`
	Method      string            `json:"method,omitempty"`
	Mode        string            `json:"mode,omitempty"`
	Priority    string            `json:"priority,omitempty"`
	Redirect    string            `json:"redirect,omitempty"`
	Referrer    string            `json:"referrer,omitempty"`
	URL         string            `json:"url,omitempty"`
}

func (f *Request) SetHeader(key, value string) {
	if f.Headers == nil {
		f.Headers = make(map[string]string)
	}
	f.Headers[key] = value
}

func (f *Request) Object() js.Value {
	var jsRequest = js.Global().Call("eval", "new Object()")
	if f.Headers != nil {
		var jsMap = js.Global().Call("eval", "new Object()")
		for key, value := range f.Headers {
			jsMap.Set(key, value)
		}
		jsRequest.Set("headers", jsMap)
	}
	if f.Body != nil {
		jsRequest.Set("body", string(f.Body))
	}
	if f.Method != "" {
		jsRequest.Set("method", strings.ToUpper(f.Method))
	}
	if f.Cache != "" {
		jsRequest.Set("cache", f.Cache)
	}
	if f.Credentials != "" {
		jsRequest.Set("credentials", f.Credentials)
	}
	if f.Destination != "" {
		jsRequest.Set("destination", f.Destination)
	}
	if f.Integrity != "" {
		jsRequest.Set("integrity", f.Integrity)
	}
	if f.Mode != "" {
		jsRequest.Set("mode", f.Mode)
	}
	if f.Priority != "" {
		jsRequest.Set("priority", f.Priority)
	}
	if f.Redirect != "" {
		jsRequest.Set("redirect", f.Redirect)
	}
	if f.Referrer != "" {
		jsRequest.Set("referrer", f.Referrer)
	}
	return jsRequest
}

type Response struct {
	Body       []byte
	Headers    map[string]string
	StatusCode int
}

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

	var fetch = js.Global().Call("fetch", options.URL, options.Object())
	if fetch.IsUndefined() {
		panic("fetch is undefined")
	}

	var respChan = make(chan *Response)
	var then = fetch.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var response = args[0]
		var headers = response.Get("headers")
		var jsHeaders = make(map[string]string)
		keys := headers.Call("keys")
		for i := 0; i < keys.Length(); i++ {
			var key = keys.Index(i).String()
			var value = headers.Call("get", key).String()
			jsHeaders[key] = value
		}
		var statusCode = response.Get("status").Int()
		var resp = &Response{
			Headers:    jsHeaders,
			StatusCode: statusCode,
		}
		response.Call("text").Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			var text = args[0].String()
			resp.Body = []byte(text)
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
