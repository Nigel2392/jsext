package state

import (
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/jsc"
)

type StatefulElement struct {
	Key        string
	Change     string
	ChangeType ChangeType
	Value      interface{}
	Elements   []SetRemover
}

// Set sets the value of the current elements included in the stateful element.
func (s *StatefulElement) Set(value any) error {
	s.Value = value
	return s.Render()
}

// Replace will replace the current elements included in the stateful element.
func (s *StatefulElement) Replace(e ...SetRemover) error {
	s.Elements = e
	return s.Render()
}

// Edit will allow you to execute a function on the stateful element.
//
// This will re-render the stateful element.
func (s *StatefulElement) Edit(fn func(*StatefulElement)) error {
	fn(s)
	return s.Render()
}

// Render the stateful elements
func (s *StatefulElement) Render() error {
	return s.renderIndex(0, len(s.Elements))
}

func (s *StatefulElement) renderIndex(start, end int) error {
	if s == nil {
		return nil
	}
	if start < 0 {
		start = 0
	}
	if end > len(s.Elements) {
		end = len(s.Elements)
	}
	if s.Elements == nil {
		return nil
	}
	if s.Change == "" {
		s.Change = "innerHTML"
	}
	if cf, ok := s.Value.(Func); ok {
		s.loopStateSetFunc(start, end, cf)
		return nil
	}
	var v, err = jsc.ValueOf(s.Value)
	if err != nil {
		return err
	}
	s.loopStateSetJS(start, end, v)
	return nil
}

func (s *StatefulElement) loopStateSetJS(start, end int, v js.Value) {
	for i := start; i < end; i++ {
		var e = s.Elements[i]
		if e == nil {
			continue
		}
		if editable, ok := e.(Editable); ok {
			editable.EditState(s.Key, s.Change, s.Value)
			continue
		}
		if s.ChangeType == ChangeTypeFunc {
			e.CallFunc(s.Change, v)
			continue
		}
		e.Set(s.Change, v)
	}
}

func (s *StatefulElement) loopStateSetFunc(start, end int, fn func() interface{}) error {
	for i := start; i < end; i++ {
		var e = s.Elements[i]
		if e == nil {
			continue
		}

		var v, err = jsc.ValueOf(fn())
		if err != nil {
			return err
		}

		if editable, ok := e.(Editable); ok {
			editable.EditState(s.Key, s.Change, v)
			continue
		}
		if s.ChangeType == ChangeTypeFunc {
			e.CallFunc(s.Change, v)
			continue
		}
		e.Set(s.Change, v)
	}
	return nil
}
