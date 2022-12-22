package router

import "strconv"

type Vars map[string]string

func (v Vars) Get(name string) string {
	return v[name]
}

func (v Vars) GetInt(name string) (int, error) {
	var str = v[name]
	return strconv.Atoi(str)
}
