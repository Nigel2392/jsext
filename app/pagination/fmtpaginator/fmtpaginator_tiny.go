//go:build js && wasm && tinygo
// +build js,wasm,tinygo

package fmtpaginator

import (
	"fmt"

	"github.com/Nigel2392/jsext/app/tokens"
)

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
	querySize   int                      `json:"-"`
	limit       int                      `json:"-"`
	fmtFetchURL string                   `json:"-"`
	fetchFunc   func([]any) ([]T, error) `json:"-"`
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
func New[T any](token *tokens.Token, formatURL string, fetchFunc func([]any) ([]T, error)) *FormatPaginator[T] {
	var p = &FormatPaginator[T]{
		Token:       token,
		fmtFetchURL: formatURL,
		querySize:   PaginatorInvalid,
		limit:       PaginatorInvalid,
		CurrentPage: 0,
		fetchFunc:   fetchFunc,
	}
	return p
}

// Return the fetched results.
func (p *FormatPaginator[T]) fetchResults(url string) ([]T, error) {
	var cli = p.client().Get(url)
	resp, err := cli.Do()
	if err != nil {
		return nil, err
	}
	jsonData, ok := resp.JSONMap()
	if jsonData == nil || !ok {
		return nil, fmt.Errorf("invalid json data")
	}
	detail, ok := jsonData["detail"]
	if ok {
		return nil, fmt.Errorf("detail: %s", detail)
	}
	count, ok := jsonData["count"]
	if !ok {
		count = 0
	}
	hasNext, ok := jsonData["has_next"]
	if !ok {
		hasNext = false
	}
	currentPage, ok := jsonData["current"]
	if !ok {
		currentPage = 0
	}
	hasPrevious, ok := jsonData["has_previous"]
	if !ok {
		hasPrevious = false
	}
	totalPages, ok := jsonData["total_pages"]
	if !ok {
		totalPages = 0
	}
	results, ok := jsonData["results"]
	if !ok {
		results = []any{}
	}
	normalizedCount, ok := count.(float64)
	if !ok {
		normalizedCount = 0
	}
	normalizedHasNext, ok := hasNext.(bool)
	if !ok {
		normalizedHasNext = false
	}
	normalizedCurrentPage, ok := currentPage.(float64)
	if !ok {
		normalizedCurrentPage = 0
	}
	normalizedHasPrevious, ok := hasPrevious.(bool)
	if !ok {
		normalizedHasPrevious = false
	}
	normalizedTotalPages, ok := totalPages.(float64)
	if !ok {
		normalizedTotalPages = 0
	}
	normalizedResults, ok := results.([]any)
	if !ok || normalizedResults == nil {
		return nil, fmt.Errorf("results is not a []any")
	}
	p.Count = int(normalizedCount)
	p.HasNext = normalizedHasNext
	p.CurrentPage = int(normalizedCurrentPage)
	p.HasPrevious = normalizedHasPrevious
	p.TotalPages = int(normalizedTotalPages)
	if p.fetchFunc == nil {
		return nil, fmt.Errorf("fetchFunc is nil")
	}
	normalized_results, err := p.fetchFunc(normalizedResults)
	if err != nil {
		return nil, err
	}
	p.Results = normalized_results
	return p.Results, nil
}
