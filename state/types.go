package state

type ChangeType int

const (
	ValueType ChangeType = iota
	FuncType  ChangeType = iota
)

// SetRemover is an interface used in state management.
//
// Optionally, types which implement SetRemover can also implement Editable.
//
// The method on Editable will be called before the Go value is converted to Javascript.
type SetRemover interface {
	Set(p string, x any)
	CallFunc(p string, x ...any)
	Remove()
}

type Editable interface {
	EditState(key string, change string, changeType ChangeType, value interface{}) error
}

type Func func() interface{}
