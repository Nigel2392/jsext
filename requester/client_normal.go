//go:build js && wasm && !tinygo
// +build js,wasm,!tinygo

package requester

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/Nigel2392/jsext/requester/encoders"
)

// APIClient is a client that can be used to make requests to a server.
func NewAPIClient() *APIClient {
	return &APIClient{
		client:        &http.Client{},
		alwaysRecover: true,
		headers:       make(map[string][]string),
	}
}

// Get a request for the specified url
func (c *APIClient) getRequest(method Methods, url string) *http.Request {
	request, err := http.NewRequest(string(method), url, nil)
	if err != nil {
		panic(err)
	}
	return request
}

// Initialize a GET request
func (c *APIClient) Get(url string) *APIClient {
	c.request = c.getRequest(GET, url)
	return c
}

// Initialize a POST request
func (c *APIClient) Post(url string) *APIClient {
	c.request = c.getRequest(POST, url)
	return c
}

// Initialize a PUT request
func (c *APIClient) Put(url string) *APIClient {
	c.request = c.getRequest(PUT, url)
	return c
}

// Initialize a PATCH request
func (c *APIClient) Patch(url string) *APIClient {
	c.request = c.getRequest(PATCH, url)
	return c
}

// Initialize a DELETE request
func (c *APIClient) Delete(url string) *APIClient {
	c.request = c.getRequest(DELETE, url)
	return c
}

// Initialize a OPTIONS request
func (c *APIClient) Options(url string) *APIClient {
	c.request = c.getRequest(OPTIONS, url)
	return c
}

// Initialize a HEAD request
func (c *APIClient) Head(url string) *APIClient {
	c.request = c.getRequest(HEAD, url)
	return c
}

// Initialize a TRACE request
func (c *APIClient) Trace(url string) *APIClient {
	c.request = c.getRequest(TRACE, url)
	return c
}

// Add form data to the request
func (c *APIClient) WithData(formData map[string]any, encoding Encoding, file ...encoders.File) *APIClient {
	if c.request == nil {
		c.errorFunc(errors.New(ErrNoRequest))
	}

	switch encoding {
	case JSON:
		c.request.Header.Set("Content-Type", "application/json")
		buf := &bytes.Buffer{}
		var err = json.NewEncoder(buf).Encode(formData)
		if err != nil {
			if c.errorFunc != nil {
				c.errorFunc(err)
			} else {
				panic(errors.New("Error encoding JSON: " + err.Error()))
			}
		}
		c.request.Body = io.NopCloser(buf)

	case FORM_URL_ENCODED:
		c.request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		var formValues = url.Values{}
		for k, v := range formData {
			var val = reflect.ValueOf(v)
			var retVal string
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			retVal = ValueToString(val)
			formValues.Add(k, retVal)
		}
		c.request.Body = io.NopCloser(bytes.NewBufferString(formValues.Encode()))

	case MULTIPART_FORM:
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		for k, v := range formData {
			var val = reflect.ValueOf(v)
			var retVal string
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			retVal = ValueToString(val)
			writer.WriteField(k, retVal)
		}
		for _, f := range file {
			part, err := writer.CreateFormFile(f.FieldName, f.FileName)
			if err != nil {
				if c.errorFunc != nil {
					c.errorFunc(err)
				} else {
					panic(errors.New("Error encoding JSON: " + err.Error()))
				}
			}
			_, err = io.Copy(part, f.Reader)
			if err != nil {
				if c.errorFunc != nil {
					c.errorFunc(err)
				} else {
					panic(errors.New("Error encoding JSON: " + err.Error()))
				}
			}
		}
		c.request.Header.Set("Content-Type", writer.FormDataContentType())
		c.request.Body = io.NopCloser(body)
	default:
		c.errorFunc(errors.New(ErrNoEncoding))
	}
	return c
}

// Make a request with url query parameters
func (c *APIClient) WithQuery(query map[string]string) *APIClient {
	if c.request == nil {
		c.errorFunc(errors.New(ErrNoRequest))
	}
	q := c.request.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	c.request.URL.RawQuery = q.Encode()
	return c
}

// Add headers to the request
func (c *APIClient) WithHeaders(headers map[string][]string) *APIClient {
	for k, v := range headers {
		c.headers[k] = append(c.headers[k], v...)
	}
	return c
}

// Add a HTTP.Cookie to the request
func (c *APIClient) WithCookie(cookie *http.Cookie) *APIClient {
	if c.request == nil {
		c.errorFunc(errors.New(ErrNoRequest))
	}
	c.request.AddCookie(cookie)
	return c
}

// Change the request before it is made
// This is useful for adding headers, cookies, etc.
func (c *APIClient) ChangeRequest(cb func(rq *http.Request)) *APIClient {
	if c.request == nil {
		c.clientErr(errors.New(ErrNoRequest))
	} else if cb == nil {
		c.clientErr(errors.New(ErrNoCallback))
	}
	cb(c.request)
	return c
}

// Set the callback function for when an error occurs
func (c *APIClient) OnError(cb func(err error) bool) *APIClient {
	c.errorFunc = cb
	return c
}

// Do not reccover when an error occurs
func (c *APIClient) NoRecover() *APIClient {
	c.alwaysRecover = false
	return c
}

func (c *APIClient) clientErr(err error) {
	if c.alwaysRecover {
		defer PrintRecover()
	}
	if err != nil {
		if c.errorFunc != nil {
			if c.errorFunc(err) {
				panic(err)
			} else {
				return
			}
		} else {
			panic(err)
		}
	}
}

// Recover from a panic and print the stack trace
func PrintRecover() any {
	if r := recover(); r != nil {
		println(string(debug.Stack()))
		println("///////////////////////////////////////////")
		println("///")
		println(r.(error).Error())
		println("///")
		println("///////////////////////////////////////////")
		return r
	}
	return nil
}

// Not used in client.
func HTTPRequest(fetchURL, method string, requestChanger func(rq *http.Request), cb func(resp *http.Response) error, onError func(err error)) { //chan struct{} {
	// var done = make(chan struct{})
	var client = http.Client{}
	var req, err = http.NewRequest(method, fetchURL, nil)
	if err != nil {
		onError(err)
		return //nil
	}
	if requestChanger != nil {
		requestChanger(req)
	}
	resp, err := client.Do(req)
	if err != nil {
		onError(err)
		return
	}
	defer resp.Body.Close()

	if cb != nil {
		err = cb(resp)
		if err != nil {
			onError(err)
			return
		}
	}
}

func ValueToString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Complex64, reflect.Complex128:
		return strconv.FormatComplex(v.Complex(), 'f', -1, 64)
	default:
		// Time
		if v.Type().String() == "time.Time" {
			return v.Interface().(time.Time).Format("2006-01-02T15:04")
		}
		return ""
	}
}
