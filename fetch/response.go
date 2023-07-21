package fetch

import (
	"syscall/js"
)

type ReadCloser interface {
	Read(p []byte) (n int, err error)
	Close() error
}

type Response struct {
	Headers    map[string][]string
	StatusCode int
	JS         js.Value
	Status     string // e.g. "200 OK"
	Request    *Request
	Body       ReadCloser
}
