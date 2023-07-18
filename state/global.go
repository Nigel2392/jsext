package state

var GlobalState = make(State)

// Set sets the state in the state map.
func Set(key string, value interface{}, change string, changeType ChangeType, e ...SetRemover) error {
	return GlobalState.Set(key, value, change, changeType, e...)
}

// Add adds the stateful element to the state map.
func Add(key string, e ...SetRemover) error {
	return GlobalState.Add(key, e...)
}

// Change changes the state in the state map.
//
// This is different from Edit, as it provides more options for changing.
func Change(key string, change string, changeType ChangeType, value interface{}) error {
	return GlobalState.Change(key, change, changeType, value)
}

// Edit changes the value of the state in the state map.
//
// This is useful for changing the value of a state without changing the change or change type.
func Edit(key string, value interface{}) error {
	return GlobalState.Edit(key, value)
}

// Delete removes the state from the state map.
//
// It also removes the elements bound to the state from the dom.
func Delete(key string, removeFromDOM bool) {
	GlobalState.Delete(key, removeFromDOM)
}

// Clear removes all the state from the state map.
//
// It also removes the elements bound to the state from the dom.
func Clear(removeFromDOM bool) {
	GlobalState.Clear(removeFromDOM)
}

// Get returns the stateful element from the state map.
func Get(key string) *StatefulElement {
	return GlobalState.Get(key)
}
