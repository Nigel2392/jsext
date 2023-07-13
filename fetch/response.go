package fetch

import "syscall/js"

type Response struct {
	Body       []byte
	Headers    map[string][]string
	StatusCode int
	JS         js.Value
}
