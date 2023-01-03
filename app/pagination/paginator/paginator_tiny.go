//go:build js && wasm && tinygo
// +build js,wasm,tinygo

package paginator

import "github.com/Nigel2392/jsext/app/tokens"

// Paginator does not work with TinyGo.
//
// This package should only be compiled using the normal Go compiler.
//
// Provides support for the token package. This is to make it easier to
// call authenticated API endpoints.
//
// (This paginator works great for Django's Rest Framework out of the box paginator.)
//
// Example of the format:
//
//	{
//	    "count": 1023,
//	    "next": "https://api.example.org/accounts/?page=5",
//	    "previous": "https://api.example.org/accounts/?page=3",
//	    "results": [
//	       …
//	    ]
//	}
//
// Simple paginator. It is used to paginate through a list of results.
//
// Fetching these results from an API endpoint.
type Paginator[T any] struct {
	Count       int                      `json:"count"`
	Next        string                   `json:"next"`
	CurrentPage string                   `json:"-"`
	Previous    string                   `json:"previous"`
	Results     []T                      `json:"results"`
	Token       *tokens.Token            `json:"-"`
	fetchURL    string                   `json:"-"`
	fetchFunc   func([]any) ([]T, error) `json:"-"`
}

// Return a new Paginator.
func New[T any](token *tokens.Token, url string, fetchFunc func([]any) ([]T, error)) *Paginator[T] {
	var p = &Paginator[T]{
		Token:     token,
		fetchURL:  url,
		fetchFunc: fetchFunc,
	}
	return p
}

// Fetch the results from the specified url.
func (p *Paginator[T]) fetchResults(url string) ([]T, error) {
	var cli = p.client().Get(url)
	resp, err := cli.Do()
	if err != nil {
		return nil, err
	}
	jsonData := resp.JSONMap()
	count, ok := jsonData["count"]
	if !ok {
		count = 0
	}
	next, ok := jsonData["next"]
	if !ok {
		next = ""
	}
	previous, ok := jsonData["previous"]
	if !ok {
		previous = ""
	}
	results, ok := jsonData["results"]
	if !ok {
		results = []any{}
	}
	p.Count = count.(int)
	p.Next = next.(string)
	p.Previous = previous.(string)
	p.CurrentPage = url
	normalized, err := p.fetchFunc(results.([]any))
	if err != nil {
		return nil, err
	}
	p.Results = normalized
	return p.Results, nil
}
