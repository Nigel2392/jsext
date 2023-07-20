package state

import (
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/errs"
)

type statefulElement[T any] struct {
	key   string
	value *js.Value
	edit  func(T, *js.Value) error
}

func (s *statefulElement[T]) Key() string {
	return s.key
}

func (s *statefulElement[T]) MarshalJS() js.Value {
	return *s.value
}

func (s *statefulElement[T]) EditState(value interface{}) error {
	v := value.(T)
	if s.edit == nil {
		return errs.Error("edit function is nil")
	}
	return s.edit(v, s.value)
}

func (s *statefulElement[T]) Remove() {
	(*s.value).Call("remove")
}

func NewElement[T any](key string, v *js.Value, edit func(T, *js.Value) error) StatefulElement {
	return &statefulElement[T]{
		key:   key,
		value: v,
		edit:  edit,
	}
}
