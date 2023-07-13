package fetch

type Response struct {
	Body       []byte
	Headers    map[string][]string
	StatusCode int
}
