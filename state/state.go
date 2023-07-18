package state

import (
	"errors"
)

func MakeSetRemoverSlice[T SetRemover](e []T) []SetRemover {
	var r = make([]SetRemover, len(e))
	for i, v := range e {
		r[i] = v
	}
	return r
}

// State is a map of stateful elements.
//
// This will manage the state of the elements.
type State map[string]*StatefulElement

// New creates a new state map.
func New(m map[string]*StatefulElement) State {
	if m == nil {
		m = make(map[string]*StatefulElement)
	}
	return m
}

// Get returns the stateful element from the state map.
func (s State) Get(key string) *StatefulElement {
	return s[key]
}

// Set sets the state in the state map.
func (s State) Set(key string, value interface{}, change string, changeType ChangeType, e ...SetRemover) error {
	if s == nil {
		s = make(map[string]*StatefulElement)
	}
	var v = &StatefulElement{
		Key:        key,
		Value:      value,
		Change:     change,
		ChangeType: changeType,
		Elements:   e,
	}
	s[key] = v
	return v.Render()
}

// Add adds the stateful element to the state map.
func (s State) Add(key string, e ...SetRemover) error {
	if s == nil {
		s = make(map[string]*StatefulElement)
	}
	if v, ok := s[key]; ok {
		var elementLen = len(v.Elements)
		v.Elements = append(v.Elements, e...)
		return v.renderIndex(elementLen, len(v.Elements))
	}
	return errors.New("state not found")
}

// Change changes the state in the state map.
//
// This is different from Edit, as it provides more options for changing.
func (s State) Change(key string, change string, changeType ChangeType, value interface{}) error {
	if s == nil {
		s = make(map[string]*StatefulElement)
	}
	if v, ok := s[key]; ok {
		v.Change = change
		v.ChangeType = changeType
		v.Value = value
		return v.Render()
	}
	return errors.New("state not found")
}

// Edit changes the value of the state in the state map.
//
// This is useful for changing the value of a state without changing the change or change type.
func (s State) Edit(key string, value interface{}) error {
	if s == nil {
		s = make(map[string]*StatefulElement)
	}
	if v, ok := s[key]; ok {
		v.Set(value)
	}
	return nil
}

// Delete removes the state from the state map.
//
// It also removes the elements bound to the state from the dom.
func (s State) Delete(key string, removeFromDOM bool) {
	if removeFromDOM {
		var v = s.Get(key)
		if v == nil {
			goto del
		}
		for _, e := range v.Elements {
			e.Remove()
		}
	}
del:
	delete(s, key)
}

// Clear removes all the state from the state map.
//
// It also removes the elements bound to the state from the dom.
func (s State) Clear(removeFromDOM bool) {
	for k := range s {
		s.Delete(k, removeFromDOM)
	}
}
