//go:build js && wasm && !tinygo
// +build js,wasm,!tinygo

package fmtpaginator

import (
	"fmt"
	"net/http"

	"github.com/Nigel2392/jsext/app/tokens"
	"github.com/Nigel2392/jsext/requester"
)

const PaginatorInvalid int = -1

// Paginator does not work with TinyGo.
//
// This package should only be compiled using the normal Go compiler.
//
// Provides support for the token package. This is to make it easier to
// call authenticated API endpoints.
//
// Example of the format:
//
//	{
//	    "has_next": true, 		// Has next page
//	    "has_previous": false, 	// Has previous page
//	    "current": 1, 		// Current page number
//	    "total_pages": 10, 		// Total number of pages
//	    "count": 5, 		// Total number of results
//	    "results": [
//	       …
//	    ]
//	}
//
// Simple paginator. It is used to paginate through a list of results.
//
// Fetching these results from an API endpoint.
type FormatPaginator[T any] struct {
	HasNext     bool `json:"has_next"`     // Has next page
	HasPrevious bool `json:"has_previous"` // Has previous page
	CurrentPage int  `json:"current"`      // Current page number
	TotalPages  int  `json:"total_pages"`  // Total pages
	Count       int  `json:"count"`
	Results     []T  `json:"results"`
	// Token is used to authenticate the request if it is not nil.
	Token *tokens.Token `json:"-"`
	// If querysize and limit are set, the url will be formatted with the page number, querysize and limit
	querySize int `json:"-"`
	limit     int `json:"-"`

	fmtFetchURL string `json:"-"`
}

// New returns a new FormatPaginator.
//
// # Format the url with the page number
//
// Example:
//
//	New(nil, "http://127.0.0.1:8000/api/users/?page=%d&querysize=%d&limit=%d")
//	New(nil, "http://127.0.0.1:8000/api/users/?page=%d&querysize=%d")
//	New(nil, "http://127.0.0.1:8000/api/users/?page=%d&limit=%d")
//	New(nil, "http://127.0.0.1:8000/api/users/?page=%d")
func New[T any](token *tokens.Token, formatURL string) *FormatPaginator[T] {
	var p = &FormatPaginator[T]{
		Token:       token,
		fmtFetchURL: formatURL,
		querySize:   PaginatorInvalid,
		limit:       PaginatorInvalid,
		CurrentPage: 0,
	}
	return p
}

// Set sets the query size and limit.
func (p *FormatPaginator[T]) Set(querySize, limit int) *FormatPaginator[T] {
	p.querySize = querySize
	p.limit = limit
	return p
}

// Format the paginator's url.
func (p *FormatPaginator[T]) url(page int) string {
	if page == PaginatorInvalid {
		page = 0
	}
	if p.querySize != PaginatorInvalid && p.limit != PaginatorInvalid {
		return fmt.Sprintf(p.fmtFetchURL, page, p.querySize, p.limit)
	} else if p.querySize != PaginatorInvalid {
		return fmt.Sprintf(p.fmtFetchURL, page, p.querySize)
	} else if p.limit != PaginatorInvalid {
		return fmt.Sprintf(p.fmtFetchURL, page, p.limit)
	}
	return fmt.Sprintf(p.fmtFetchURL, page)
}

// Next returns the next page number.
func (p *FormatPaginator[T]) Next() int {
	return p.CurrentPage + 1
}

// Previous returns the previous page number.
func (p *FormatPaginator[T]) Previous() int {
	return p.CurrentPage - 1
}

// FetchNext fetches the next page of results.
func (p *FormatPaginator[T]) FetchNext() ([]T, error) {
	if !p.HasNext {
		return nil, nil
	}
	p.CurrentPage++
	return p.fetchResults(p.url(p.CurrentPage))
}

// Fetch fetches the current page of results.
func (p *FormatPaginator[T]) Fetch(page ...int) ([]T, error) {
	if len(page) > 0 {
		return p.fetchResults(p.url(page[0]))
	}
	if p.CurrentPage > 1 {
		return p.FetchNext()
	}
	if p.CurrentPage < 1 {
		p.CurrentPage = 1
	}
	return p.fetchResults(p.url(p.CurrentPage))
}

// FetchPrevious fetches the previous page of results.
func (p *FormatPaginator[T]) FetchPrevious() ([]T, error) {
	if !p.HasPrevious {
		return nil, nil
	}
	p.CurrentPage--
	return p.fetchResults(p.url(p.CurrentPage))
}

// Return the fetched results.
func (p *FormatPaginator[T]) fetchResults(url string) ([]T, error) {
	var cli = p.client().Get(url)
	var newP = &FormatPaginator[T]{}
	var waiter = make(chan *FormatPaginator[T])
	cli.DoDecodeTo(newP, requester.JSON, func(resp *http.Response, strct interface{}) {
		p.Count = newP.Count
		p.HasNext = newP.HasNext
		p.CurrentPage = newP.CurrentPage
		p.HasPrevious = newP.HasPrevious
		p.TotalPages = newP.TotalPages
		p.Results = newP.Results
		waiter <- newP
	})
	return (<-waiter).Results, nil
}

// Default client.
func (p *FormatPaginator[T]) client() *requester.APIClient {
	if p.Token != nil && p.Token.AccessToken != "" {
		return p.Token.Client()
	}
	return requester.NewAPIClient()
}
