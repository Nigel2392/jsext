//go:build js && wasm && !tinygo
// +build js,wasm,!tinygo

package requester

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
)

// APIClient is a client that can be used to execute http requests.
// - Can be used to execute GET, POST, PUT, DELETE, PATCH requests.
type APIClient struct {
	// Client is the http client that will be used to execute the request.
	client *http.Client
	// Request is the request that will be executed.
	request *http.Request
	// Error func is used to handle errors that occur during the request.
	// - Can be specified, but is not required.
	// - If not set will panic.
	// - If set and returns true will panic
	// - Always recovers from panics when set
	errorFunc     func(err error) bool
	before        func()
	after         func()
	alwaysRecover bool
}

// Function to execute before the request is executed
func (c *APIClient) Before(cb func()) {
	c.before = cb
}

// Function to execute after the request is executed
func (c *APIClient) After(cb func()) {
	c.after = cb
}

// Execute the request -> APIClient.exec
func (c *APIClient) Do(cb func(resp *http.Response)) {
	if c.request == nil {
		c.errorFunc(errors.New(ErrNoRequest))
	} else if cb == nil {
		c.errorFunc(errors.New(ErrNoCallback))
	}
	go func() {
		if c.before != nil {
			c.before()
		}
		c.exec(cb)
		if c.after != nil {
			c.after()
		}
	}()
}

// Make a request and return the response body decoded into the specified parameter. -> APIClient.Do  -> APIClient.exec
func (c *APIClient) DoStruct(decodeTo interface{}, encoding Encoding, cb func(resp *http.Response, strct interface{})) {
	var newCallback = func(resp *http.Response) {
		switch encoding {
		case JSON:
			if err := json.NewDecoder(resp.Body).Decode(decodeTo); err != nil {
				c.clientErr(err)
			}
		case XML:
			if err := xml.NewDecoder(resp.Body).Decode(decodeTo); err != nil {
				c.clientErr(err)
			}
		default:
			c.clientErr(errors.New(ErrNoEncoding))
		}

		if cb != nil {
			cb(resp, decodeTo)
		}
	}
	c.Do(newCallback)
}

// Execute the request
func (c *APIClient) exec(cb func(resp *http.Response)) error {
	var resp, err = c.client.Do(c.request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	cb(resp)
	return nil
}
