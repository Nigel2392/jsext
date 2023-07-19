package stateful

import (
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/errs"
	"github.com/Nigel2392/jsext/v2/jsc"
	"github.com/Nigel2392/jsext/v2/state"
)

type Elements[T any] struct {
	Elements []js.Value
	Data     T
	EditFunc func(T, js.Value) error
}

func NewElements[T any](s ...js.Value) *Elements[T] {
	if s == nil {
		s = make([]js.Value, 0)
	}
	return &Elements[T]{
		Elements: s,
	}
}

func (s *Elements[T]) Set(name string, value any) {
	if s.Elements == nil {
		return
	}
	for _, elem := range s.Elements {
		elem.Set(name, value)
	}
}

func (s *Elements[T]) CallFunc(funcName string, args ...any) {
	if s == nil {
		return
	}
	for _, elem := range s.Elements {
		elem.Call(funcName, args...)
	}
}

func (s *Elements[T]) Remove() {
	if s == nil {
		return
	}
	for _, elem := range s.Elements {
		elem.Call("remove")
	}
}

func (s *Elements[T]) EditState(key string, change string, changeType state.ChangeType, value interface{}) error {
	if s == nil {
		return nil
	}
	var err error
	var v = value.(T)
	for _, elem := range s.Elements {
		if s.EditFunc != nil {
			err = s.EditFunc(v, elem)
			if err != nil {
				return err
			}
			continue
		}
		var v, err = jsc.ValueOf(value)
		if err != nil {
			return err
		}
		switch {
		case changeType == state.ValueType:
			elem.Set(change, v)
		case changeType == state.FuncType:
			elem.Call(change, v)
		default:
			return errs.Error("invalid change type")
		}
	}
	return nil
}

func (s *Elements[T]) AppendChild(e ...js.Value) {
	if s == nil {
		return
	}
	if s == nil {
		s.Elements = make([]js.Value, 0)
	}
	s.Elements = append(s.Elements, e...)
}
