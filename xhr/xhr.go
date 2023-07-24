package xhr

import (
	"syscall/js"
	"time"

	"github.com/Nigel2392/jsext/v2/encoding"
	"github.com/Nigel2392/jsext/v2/errs"
)

var xhrGlobal = js.Global().Get("XMLHttpRequest")

type XMLHttpRequest struct {
	Timeout       time.Duration
	Upload        *XMLHttpRequestUpload
	ResponseType  string
	StatusCodeMap map[int]func(*XMLHttpRequest)
	*js.Value
}

func New() *XMLHttpRequest {
	var xhr = xhrGlobal.New()
	var upload = xhr.Get("upload")
	var x = &XMLHttpRequest{
		Value: &xhr,
		Upload: &XMLHttpRequestUpload{
			Value: &upload,
		},
	}
	return x
}

func OpenNew(method, url string, user User) *XMLHttpRequest {
	var x = New()
	x.Open(method, url, user)
	return x
}

func (x *XMLHttpRequest) Open(method, url string, user User) {
	var (
		username string
		password string
	)

	if user == nil {
		x.Call("open", method, url, true)
		return
	}

	username = user.Username()
	password = user.Password()

	if username != "" && password != "" {
		x.Call("open", method, url, true, username, password)
		return
	}

	if username != "" {
		x.Call("open", method, url, true, username)
		return
	}

	x.Call("open", method, url, true)
}

func (x *XMLHttpRequest) MarshalJS() js.Value {
	return *x.Value
}

func (x *XMLHttpRequest) SetHeader(key, value string) {
	x.Call("setRequestHeader", key, value)
}

func (x *XMLHttpRequest) GetHeader(key string) string {
	return x.Call("getResponseHeader", key).String()
}

func (x *XMLHttpRequest) Status() int {
	return x.Get("status").Int()
}

func (x *XMLHttpRequest) StatusText() string {
	return x.Get("statusText").String()
}

func (x *XMLHttpRequest) Response() js.Value {
	return x.Get("response")
}

func (x *XMLHttpRequest) IsDone() bool {
	return x.Get("readyState").Int() == 4
}

const (
	STATE_UNSENT = iota
	STATE_OPENED
	STATE_HEADERS_RECEIVED
	STATE_LOADING
	STATE_DONE
)

func (x *XMLHttpRequest) OnStatus(statusCode int, f func(*XMLHttpRequest)) {
	if x.StatusCodeMap == nil {
		x.StatusCodeMap = make(map[int]func(*XMLHttpRequest))
	}
	x.StatusCodeMap[statusCode] = f
}

func (x *XMLHttpRequest) OnReadyStateChange(f func(state int, xhr *XMLHttpRequest)) js.Func {
	if f == nil {
		x.Set("onreadystatechange", nil)
		return js.Func{}
	}
	var fn = js.FuncOf(func(this js.Value, args []js.Value) any {
		var (
			readyState = x.Get("readyState").Int()
			status     = x.Get("status").Int()
		)
		f(readyState, x)
		if readyState == STATE_DONE {
			if f, ok := x.StatusCodeMap[status]; ok {
				f(x)
			}
		}
		return nil
	})
	x.Set("onreadystatechange", fn)
	return fn
}

const (
	ErrAborted = errs.Error("XHR aborted")
	ErrTimeout = errs.Error("XHR timeout")
)

func (x *XMLHttpRequest) Send(data any) (response js.Value, err error) {

	if x.Timeout <= 0 {
		x.Timeout = 5 * time.Second
	}

	if x.ResponseType != "" {
		x.Set("responseType", x.ResponseType)
	}

	var (
		onLoad   js.Func
		onErr    js.Func
		onAbort  js.Func
		closed   = make(chan struct{})
		respChan = make(chan js.Value)
		errChan  = make(chan error)
	)

	onLoad = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		select {
		case <-closed:
			return nil
		default:
			respChan <- x.Response()
		}
		return nil
	})
	onErr = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			var event = args[0]
			select {
			case <-closed:
				return nil
			default:
				var textContent = event.Get("textContent")
				if str := textContent.String(); textContent.Type() == js.TypeString && str != "" {
					errChan <- errs.Error(event.Get("textContent").String())
				} else {
					errChan <- errs.Error("XHR error occurred")
				}
			}
		}
		return nil
	})
	onAbort = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		select {
		case <-closed:
			return nil
		default:
			errChan <- ErrAborted
		}
		return nil
	})

	defer func() {
		close(closed)
		close(respChan)
		close(errChan)

		x.Set("onload", js.Null())
		x.Set("onerror", js.Null())
		x.Set("onabort", js.Null())

		onLoad.Release()
		onErr.Release()
		onAbort.Release()
	}()

	x.Set("onload", onLoad)
	x.Set("onerror", onErr)
	x.Set("onabort", onAbort)

	if data == nil {
		x.Call("send")
		goto wait
	}

	switch data := data.(type) {
	case string:
		x.Call("send", data)
	case []byte:
		var buf = js.Global().Get("Uint8Array").New(len(data))
		js.CopyBytesToJS(buf, data)
		x.Call("send", buf)
	case js.Value:
		if data.Type() == js.TypeString {
			x.Call("send", data.String())
			break
		}
		if data.InstanceOf(js.Global().Get("ArrayBuffer")) ||
			data.InstanceOf(js.Global().Get("Int8Array")) ||
			data.InstanceOf(js.Global().Get("Blob")) ||
			data.InstanceOf(js.Global().Get("Document")) ||
			data.InstanceOf(js.Global().Get("FormData")) {
			x.Call("send", data)
			break
		}
	default:
		data, err = encoding.EncodeJSON[string](data) // string, error
		if err != nil {
			return js.Null(), err
		}
		x.Call("send", data)
	}

wait:
	select {
	case <-time.After(x.Timeout):
		return js.Null(), ErrTimeout
	case response = <-respChan:
	case err = <-errChan:
		return js.Null(), err
	}
	return response, nil
}
