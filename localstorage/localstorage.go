package localstorage

import (
	"errors"
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/jsc"
)

var localStorage = js.Global().Get("localStorage")

// Set a key value pair in localStorage
func Set(key, value string) error {
	if localStorage.IsUndefined() {
		return errors.New("localStorage is undefined")
	}
	localStorage.Call("setItem", key, value)
	return nil
}

// Get a value from localStorage
func Get(key string) (string, error) {
	if localStorage.IsUndefined() {
		return "", errors.New("localStorage is undefined")
	}
	var item = localStorage.Call("getItem", key)
	if item.IsNull() || item.IsUndefined() {
		return "", errors.New("key not found")
	}
	return item.String(), nil
}

// Remove a key value pair from localStorage
func Remove(key string) error {
	if localStorage.IsUndefined() {
		return errors.New("localStorage is undefined")
	}
	localStorage.Call("removeItem", key)
	return nil
}

// Clear all key value pairs from localStorage
func Clear() error {
	if localStorage.IsUndefined() {
		return errors.New("localStorage is undefined")
	}
	localStorage.Call("clear")
	return nil
}

var json = js.Global().Get("JSON")

// Try to set any object by first converting it to a js.Value,
// then converting it to a string, shelling out to JSON.stringify,
func UnsafeSet(key, value interface{}) error {
	if localStorage.IsUndefined() {
		return errors.New("localStorage is undefined")
	}
	var v, err = jsc.ValueOf(value)
	if err != nil {
		return err
	}
	var s = json.Call("stringify", v)
	localStorage.Call("setItem", key, s)
	return nil
}

// Try to get any object by first getting the string from localStorage,
// then shell out to JSON.parse and scan the object into dst.
func UnsafeGet(key string, dst interface{}) error {
	if localStorage.IsUndefined() {
		return errors.New("localStorage is undefined")
	}
	var item = localStorage.Call("getItem", key)
	if item.IsNull() || item.IsUndefined() {
		return errors.New("key not found")
	}
	var v = json.Call("parse", item)
	return jsc.Scan(v, dst)
}
