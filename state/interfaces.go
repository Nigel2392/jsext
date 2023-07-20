package state

import "syscall/js"

type Keyer interface {
	Key() string
}

type Key string

func (k Key) Key() string {
	return string(k)
}

type Editor interface {
	EditState(value interface{}) error
}

type StatefulElement interface {
	Keyer
	Editor
	MarshalJS() js.Value
	Remove()
}
