package fetch

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"syscall/js"
)

type Request struct {
	Body        []byte
	GetBody     func() (io.ReadCloser, error)
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
	ctx         context.Context
}

func NewRequest(method, url string) *Request {
	return &Request{
		URL: url,
		ctx: context.Background(),
	}
}

func (f *Request) Context() context.Context {
	if f.ctx == nil {
		f.ctx = context.Background()
	}
	return f.ctx
}

func (f *Request) SetContext(ctx context.Context) {
	f.ctx = ctx
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
	switch body := body.(type) {
	case []byte:
		f.Body = body
	case string:
		f.Body = []byte(body)
	case io.ReadCloser:
		var buf bytes.Buffer
		if _, err = io.Copy(&buf, body); err != nil {
			return err
		}
		f.Body = buf.Bytes()
	case io.Reader:
		var buf bytes.Buffer
		if _, err = io.Copy(&buf, body); err != nil {
			return err
		}
		f.Body = buf.Bytes()
	case func() (io.ReadCloser, error):
		f.GetBody = body
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

func (f *Request) MarshalJS() (js.Value, error) {
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
	} else if f.GetBody != nil {
		var reader, err = f.GetBody()
		if err != nil {
			return js.Null(), err
		}
		var buf bytes.Buffer
		if _, err = io.Copy(&buf, reader); err != nil {
			return js.Null(), err
		}
		var jsBody = js.Global().Get("Uint8Array").New(len(buf.Bytes()))
		js.CopyBytesToJS(jsBody, buf.Bytes())
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
	return jsRequest, nil
}
