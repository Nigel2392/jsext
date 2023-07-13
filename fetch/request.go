package fetch

import (
	"bytes"
	"encoding/json"
	"io"
	"syscall/js"
)

type Request struct {
	Body        []byte
	Cache       string
	Credentials string
	Destination string
	Headers     map[string][]string
	Integrity   string
	Method      string
	Mode        string
	Priority    string
	Redirect    string
	Referrer    string
	URL         string
}

func (f *Request) SetHeader(key, value string) {
	if f.Headers == nil {
		f.Headers = make(map[string][]string)
	}
	f.Headers[key] = []string{value}
}

func (f *Request) AddHeader(key, value string) {
	if f.Headers == nil {
		f.Headers = make(map[string][]string)
	}
	f.Headers[key] = append(f.Headers[key], value)
}

func (f *Request) DeleteHeader(key string) {
	if f.Headers == nil {
		f.Headers = make(map[string][]string)
	}
	delete(f.Headers, key)
}

func (f *Request) SetBody(body any) (err error) {
	switch body.(type) {
	case []byte:
		f.Body = body.([]byte)
	case string:
		f.Body = []byte(body.(string))
	case io.Reader:
		var buf bytes.Buffer
		if _, err = io.Copy(&buf, body.(io.Reader)); err != nil {
			return err
		}
		f.Body = buf.Bytes()
	case io.ReadCloser:
		var buf bytes.Buffer
		if _, err = io.Copy(&buf, body.(io.ReadCloser)); err != nil {
			return err
		}
		f.Body = buf.Bytes()
	default:
		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return err
		}
		f.Body = buf.Bytes()
		f.SetHeader("Content-Type", "application/json")
	}
	return nil
}

func (f *Request) MarshalJS() js.Value {
	var jsRequest = js.Global().Get("Object").New()
	if f.Headers != nil {
		var jsMap = js.Global().Get("Object").New()
		for key, value := range f.Headers {
			var jsValue = js.Global().Get("Array").New()
			for _, v := range value {
				jsValue.Call("push", v)
			}
			jsMap.Set(key, jsValue)
		}
		jsRequest.Set("headers", jsMap)
	}
	if f.Body != nil {
		var jsBody = js.Global().Get("Uint8Array").New(len(f.Body))
		js.CopyBytesToJS(jsBody, f.Body)
		jsRequest.Set("body", jsBody)
	}
	if f.Method != "" {
		jsRequest.Set("method", f.Method)
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
