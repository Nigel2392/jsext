//go:build js && wasm
// +build js,wasm

package fmtpaginator

import (
	"fmt"

	"github.com/Nigel2392/jsext/requester"
)

const PaginatorInvalid int = -1

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

// Default client.
func (p *FormatPaginator[T]) client() *requester.APIClient {
	if p.Token != nil && p.Token.AccessToken != "" {
		return p.Token.Client()
	}
	return requester.NewAPIClient()
}
