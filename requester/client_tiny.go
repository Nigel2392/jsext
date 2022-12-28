//go:build js && wasm && tinygo
// +build js,wasm,tinygo

package requester

import (
	"github.com/Nigel2392/jsext/requester/fetch"
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
func (c *APIClient) WithData(formData map[string]interface{}, encoding Encoding, file ...File) *APIClient {
	switch encoding {
	case JSON:
		c.Request.Headers["Content-Type"] = "application/json"
		c.Request.Body = fetch.MarshalMap(formData)
	case FORM_URL_ENCODED:
		panic(FORM_URL_ENCODED + " is not supported yet!")
	case MULTIPART_FORM:
		panic(MULTIPART_FORM + " is not supported yet!")
	}
	return c
}

// Send the request
func (c *APIClient) Do() *fetch.Response {
	if c.before != nil {
		c.before()
	}
	var resp = fetch.Fetch(c.Request)
	if c.after != nil {
		c.after()
	}
	return resp
}

// Function to execute before the request is executed
func (c *APIClient) Before(cb func()) {
	c.before = cb
}

// Function to execute after the request is executed
func (c *APIClient) After(cb func()) {
	c.after = cb
}
