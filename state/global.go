package state

var GlobalState = make(State)

func Set(key string, value interface{}, change string, changeType ChangeType, e ...SetRemover) error {
	return GlobalState.Set(key, value, change, changeType, e...)
}

func Add(key string, e ...SetRemover) error {
	return GlobalState.Add(key, e...)
}

func Change(key string, change string, changeType ChangeType, value interface{}) error {
	return GlobalState.Change(key, change, changeType, value)
}

func Edit(key string, value interface{}) error {
	return GlobalState.Edit(key, value)
}

func Delete(key string, removeFromDOM bool) {
	GlobalState.Delete(key, removeFromDOM)
}

func Clear(removeFromDOM bool) {
	GlobalState.Clear(removeFromDOM)
}

func Get(key string) *StatefulElement {
	return GlobalState.Get(key)
}
