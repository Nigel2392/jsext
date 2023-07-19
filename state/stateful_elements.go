package state

import (
	"github.com/Nigel2392/jsext/v2/errs"
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
func (s *StatefulElement) Remove(e ...SetRemover) error {
	for _, v := range e {
	inner:
		for i, e := range s.Elements {
			if e == v {
				s.removeIndex(i)
				e.Remove()
				break inner
			}
		}
	}
	return s.Render()
}

// removeIndex will remove the current elements included in the stateful element.
func (s *StatefulElement) removeIndex(i int) {
	if i < 0 || i >= len(s.Elements) {
		return
	}
	if len(s.Elements) == 1 {
		s.Elements = make([]SetRemover, 0)
		return
	}
	if i == 0 {
		s.Elements = s.Elements[1:]
		return
	}
	if i == len(s.Elements)-1 {
		s.Elements = s.Elements[:len(s.Elements)-1]
		return
	}
	s.Elements = append(s.Elements[:i], s.Elements[i+1:]...)
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
		return s.loopStateSetFunc(start, end, cf)
	}
	var v, err = jsc.ValueOf(s.Value)
	if err != nil {
		return err
	}
	var cf = func() interface{} {
		return v
	}
	return s.loopStateSetFunc(start, end, cf)
}

func (s *StatefulElement) loopStateSetFunc(start, end int, fn func() interface{}) error {
	var (
		i   int
		v   any
		e   SetRemover
		err error
	)
	for i = start; i < end; i++ {
		e = s.Elements[i]
		if e == nil {
			continue
		}
		v = fn()
		if editable, ok := e.(Editable); ok {
			err = editable.EditState(s.Key, s.Change, s.ChangeType, v)
			if err != nil {
				return err
			}
		}
		v, err = jsc.ValueOf(v)
		if err != nil {
			return err
		}
		switch {
		case s.ChangeType == ValueType:
			e.Set(s.Change, v)
		case s.ChangeType == FuncType:
			e.CallFunc(s.Change, v)
		default:
			return errs.Error("invalid change type")
		}
	}
	return nil
}
