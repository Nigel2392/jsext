//go:build js && wasm && tinygo
// +build js,wasm,tinygo

package requester

import (
	"errors"

	"github.com/Nigel2392/jsext/framework/requester/encoders"
	"github.com/Nigel2392/jsext/framework/requester/fetch"
)

// This client works a little bit differently from the normal APIClient.
// This is due to not everything being supported in TinyGo,
// Or being too much of a hassle to implement.
type APIClient struct {
	Request *fetch.Request
	before  func()
	after   func()
}

func NewAPIClient() *APIClient {
	return &APIClient{
		Request: &fetch.Request{
			Headers: make(map[string]string),
		},
	}
}

func (c *APIClient) Get(url string) *APIClient {
	c.Request.Method = "GET"
	c.Request.URL = url
	return c
}

func (c *APIClient) Post(url string) *APIClient {
	c.Request.Method = "POST"
	c.Request.URL = url
	return c
}

func (c *APIClient) Put(url string) *APIClient {
	c.Request.Method = "PUT"
	c.Request.URL = url
	return c
}

func (c *APIClient) Patch(url string) *APIClient {
	c.Request.Method = "PATCH"
	c.Request.URL = url
	return c
}

func (c *APIClient) Delete(url string) *APIClient {
	c.Request.Method = "DELETE"
	c.Request.URL = url
	return c
}

func (c *APIClient) Head(url string) *APIClient {
	c.Request.Method = "HEAD"
	c.Request.URL = url
	return c
}

// Add headers to the request
func (c *APIClient) WithHeaders(headers map[string]string) *APIClient {
	if c.Request.Headers == nil {
		c.Request.Headers = make(map[string]string)
	}
	for k, v := range headers {
		c.Request.Headers[k] = v
	}
	return c
}

// Add data to the request
func (c *APIClient) WithData(formData map[string]interface{}, encoding Encoding, file ...encoders.File) *APIClient {
	switch encoding {
	case JSON:
		c.Request.Headers["Content-Type"] = "application/json"
		c.Request.Body = encoders.MarshalMap(formData)
	case FORM_URL_ENCODED:
		c.Request.Headers["Content-Type"] = "application/x-www-form-urlencoded"
		c.Request.Body = encoders.ToURLValues(formData)
	case MULTIPART_FORM:
		c.Request.Headers["Content-Type"] = "multipart/form-data"
		c.Request.Body = encoders.ToMultipart(formData, file...)
	}
	return c
}

// Send the request
func (c *APIClient) Do() (*fetch.Response, error) {
	if c.before != nil {
		c.before()
	}
	var resp, err = fetch.Fetch(c.Request)
	if err != nil {
		return nil, err
	}
	if c.after != nil {
		c.after()
	}
	return resp, nil
}

func (c *APIClient) DoDecode(encType Encoding) (map[string]interface{}, *fetch.Response, error) {
	var resp, err = c.Do()
	if err != nil {
		return nil, nil, err
	}
	var data map[string]interface{}
	switch encType {
	case JSON:
		var ok bool
		data, ok = resp.JSONMap()
		if !ok {
			return nil, nil, errors.New("could not decode json")
		}
	case FORM_URL_ENCODED:
		panic("Form url encoded is not supported yet!")
	case MULTIPART_FORM:
		panic("Multipart form is not supported yet!")
	}
	return data, resp, nil
}

// Function to execute before the request is executed
func (c *APIClient) Before(cb func()) {
	c.before = cb
}

// Function to execute after the request is executed
func (c *APIClient) After(cb func()) {
	c.after = cb
}
