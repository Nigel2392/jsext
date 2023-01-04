package vars

import "strconv"

type Vars map[string]string

// Get a variable from the vars map
func (v Vars) Get(name string) string {
	return v[name]
}

// Get an integer from the vars map
func (v Vars) GetInt(name string) (int, error) {
	var str = v[name]
	return strconv.Atoi(str)
}
