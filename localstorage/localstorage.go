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
	return localStorage.Call("getItem", key).String(), nil
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

func UnsafeSet(key, value interface{}) error {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
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

func UnsafeGet(key string, dst interface{}) error {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	if localStorage.IsUndefined() {
		return errors.New("localStorage is undefined")
	}
	var s = localStorage.Call("getItem", key).String()
	var v = json.Call("parse", s)
	return jsc.Scan(v, dst)
}
