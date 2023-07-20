package state

import (
	"fmt"
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/errs"
)

type ListStateFlags uint32

const (
	F_APPEND ListStateFlags = 1 << iota
	F_PREPEND
	F_NONE
)

type State struct {
	Root     js.Value
	Elements map[string]StatefulElement
	Flags    ListStateFlags
}

func New(root js.Value) *State {
	return &State{
		Root:     root,
		Elements: make(map[string]StatefulElement),
	}
}

func (s *State) With(value interface{}, e StatefulElement) error {
	if s == nil {
		return errs.Error("stateful is nil")
	}
	if s.Elements == nil {
		s.Elements = make(map[string]StatefulElement)
	}
	return s.updateOrAdd(e, value)
}

func (s *State) Update(key Keyer, v interface{}) error {
	if s == nil {
		return errs.Error("stateful is nil")
	}
	var e, ok = s.Elements[key.Key()]
	if !ok {
		var keys = make([]string, 0, len(s.Elements))
		for k := range s.Elements {
			keys = append(keys, k)
		}
		fmt.Println("Keys", keys, "Key", key.Key())
		return errs.Error("key not found")
	}
	return e.EditState(v)
}

func (s *State) Without(e ...Keyer) error {
	if s == nil {
		return errs.Error("stateful is nil")
	}
	for _, elem := range e {
		var key = elem.Key()
		var e, ok = s.Elements[key]
		if !ok {
			continue
		}
		e.Remove()
		delete(s.Elements, key)
	}
	return nil
}

func (s *State) Clear() {
	if s == nil {
		return
	}
	for _, elem := range s.Elements {
		elem.Remove()
	}
	s.Elements = make(map[string]StatefulElement)
}

func (s *State) updateOrAdd(e StatefulElement, v interface{}) error {
	if e == nil {
		return errs.Error("stateful element is nil")
	}
	if _, ok := s.Elements[e.Key()]; ok {
		return s.update(e, v)
	}
	return s.add(e, v)
}

func (s *State) update(e StatefulElement, v interface{}) error {
	var err error
	var oldKey = e.Key()
	err = e.EditState(v)
	if err != nil {
		return err
	}
	var newKey = e.Key()
	if oldKey != newKey {
		delete(s.Elements, oldKey)
	}
	fmt.Println("Updating key", newKey)
	s.Elements[newKey] = e
	return nil
}

func (s *State) add(e StatefulElement, v interface{}) error {
	fmt.Println("Adding key", e.Key())
	if err := s.update(e, v); err != nil {
		return err
	}
	var keys = make([]string, 0, len(s.Elements))
	for k := range s.Elements {
		keys = append(keys, k)
	}
	fmt.Println("Keys", keys)
	fmt.Println("Added key", e.Key())

	if !s.Root.IsNull() && !s.Root.IsUndefined() {
		switch {
		case s.Flags&F_APPEND != 0:
			s.Root.Call("appendChild", e.MarshalJS())
		case s.Flags&F_PREPEND != 0:
			s.Root.Call("prepend", e.MarshalJS())
		default:
			return nil
		}
	}
	return nil
}
