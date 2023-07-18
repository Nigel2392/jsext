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

type State map[string]*StatefulElement

func (s State) Get(key string) *StatefulElement {
	return s[key]
}

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
