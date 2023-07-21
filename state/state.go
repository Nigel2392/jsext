package state

import (
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/errs"
)

type StateFlags uint32

const (
	F_APPEND StateFlags = 1 << iota
	F_PREPEND
	F_NONE
)

type State struct {
	Root     js.Value
	Elements map[string][]StatefulElement
	Flags    StateFlags
	OnUpdate func()
}

func New(root js.Value) *State {
	return &State{
		Root:     root,
		Elements: make(map[string][]StatefulElement),
	}
}

func (s *State) With(value interface{}, e StatefulElement) error {
	if s == nil {
		return errs.Error("stateful is nil")
	}
	if s.Elements == nil {
		s.Elements = make(map[string][]StatefulElement)
	}
	return s.updateOrAdd(e, value)
}

func (s *State) Update(key Keyer, v interface{}) error {
	if s == nil {
		return errs.Error("stateful is nil")
	}
	var e, ok = s.Elements[key.Key()]
	if !ok {
		return errs.Error("key not found")
	}
	var err error
	for _, elem := range e {
		err = elem.EditState(v)
		if err != nil {
			return err
		}
	}
	if s.OnUpdate != nil {
		s.OnUpdate()
	}
	return nil
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
		delete(s.Elements, key)
		if s.Root.IsNull() || s.Root.IsUndefined() {
			continue
		}
		for _, elem := range e {
			elem.Remove()
		}
	}
	if s.OnUpdate != nil {
		s.OnUpdate()
	}
	return nil
}

func (s *State) Clear() {
	if s == nil {
		return
	}
	for _, elem := range s.Elements {
		for _, e := range elem {
			if s.Root.IsNull() || s.Root.IsUndefined() {
				continue
			}
			e.Remove()
		}
	}
	s.Elements = make(map[string][]StatefulElement)
	if s.OnUpdate != nil {
		s.OnUpdate()
	}
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
	if s.OnUpdate != nil {
		s.OnUpdate()
	}
	var newKey = e.Key()
	if oldKey != newKey {
		delete(s.Elements, oldKey)
	}
	s.Elements[newKey] = append(s.Elements[newKey], e)
	return nil
}

func (s *State) add(e StatefulElement, v interface{}) error {
	if err := s.update(e, v); err != nil {
		return err
	}

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
