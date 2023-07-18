package state

type ChangeType int

const (
	ChangeTypeValue ChangeType = iota
	ChangeTypeFunc  ChangeType = iota
)

type SetRemover interface {
	Set(p string, x any)
	CallFunc(p string, x ...any)
	Remove()
}

type Editable interface {
	EditState(key string, change string, value interface{}) error
}

type Func func() interface{}
