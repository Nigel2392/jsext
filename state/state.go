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

type State interface {
	Get(key string) *StatefulElement
	Set(key string, value interface{}, change string, changeType ChangeType, e ...SetRemover) error
	Add(key string, e ...SetRemover) error
	Change(key string, change string, changeType ChangeType, value interface{}) error
	Replace(key string, e ...SetRemover) error
	Edit(key string, value interface{}) error
	Delete(key string, removeFromDOM bool)
	Clear(removeFromDOM bool)
}

// State is a map of stateful elements.
//
// This will manage the state of the elements.
type state map[string]*StatefulElement

// New creates a new state map.
func New(m map[string]*StatefulElement) State {
	if m == nil {
		m = make(map[string]*StatefulElement)
	}
	return state(m)
}

// Get returns the stateful element from the state map.
func (s state) Get(key string) *StatefulElement {
	return s[key]
}

// Set sets the state in the state map.
func (s state) Set(key string, value interface{}, change string, changeType ChangeType, e ...SetRemover) error {
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
func (s state) Add(key string, e ...SetRemover) error {
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
func (s state) Change(key string, change string, changeType ChangeType, value interface{}) error {
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

// Replace will replace the current elements included in the stateful element.
func (s state) Replace(key string, e ...SetRemover) error {
	if s == nil {
		s = make(map[string]*StatefulElement)
	}
	if v, ok := s[key]; ok {
		v.Replace(e...)
	}
	return errors.New("state not found")
}

// Edit changes the value of the state in the state map.
//
// This is useful for changing the value of a state without changing the change or change type.
func (s state) Edit(key string, value interface{}) error {
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
func (s state) Delete(key string, removeFromDOM bool) {
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
func (s state) Clear(removeFromDOM bool) {
	for k := range s {
		s.Delete(k, removeFromDOM)
	}
}
