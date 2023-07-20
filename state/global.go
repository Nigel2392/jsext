package state

import "syscall/js"

var GlobalState = New(js.Null())

// Get returns the stateful element from the state map.
func Get(key Keyer) StatefulElement {
	return GlobalState.Elements[key.Key()]
}

// With adds the stateful element to the state map.
func With(value interface{}, e StatefulElement) error {
	return GlobalState.With(value, e)
}

// Without removes the stateful element from the state map.
func Without(e ...Keyer) error {
	return GlobalState.Without(e...)
}

// Update updates the stateful element in the state map.
func Update(key Keyer, v interface{}) error {
	return GlobalState.Update(key, v)
}
