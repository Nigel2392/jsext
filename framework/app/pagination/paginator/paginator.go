//go:build js && wasm && !tinygo
// +build js,wasm,!tinygo

package paginator

import (
	"net/http"

	"github.com/Nigel2392/jsext/framework/app/tokens"
	"github.com/Nigel2392/jsext/framework/requester"
)

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
	Count       int           `json:"count"`
	Next        string        `json:"next"`
	CurrentPage string        `json:"-"`
	Previous    string        `json:"previous"`
	Results     []T           `json:"results"`
	Token       *tokens.Token `json:"-"`
	fetchURL    string        `json:"-"`
}

// Return a new Paginator.
func New[T any](token *tokens.Token, url string) *Paginator[T] {
	var p = &Paginator[T]{
		Token:    token,
		fetchURL: url,
	}
	return p
}

// Fetch the results from the specified url.
func (p *Paginator[T]) fetchResults(url string) ([]T, error) {
	var cli = p.client().Get(url)
	var newP = &Paginator[T]{}
	var waiter = make(chan *Paginator[T])
	cli.DoDecodeTo(newP, requester.JSON, func(resp *http.Response, strct interface{}) {
		p.Count = newP.Count
		p.Next = newP.Next
		p.CurrentPage = newP.CurrentPage
		p.Previous = newP.Previous
		p.Results = newP.Results
		waiter <- newP
	})
	return (<-waiter).Results, nil
}
