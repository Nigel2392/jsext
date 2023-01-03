package paginator

import "github.com/Nigel2392/jsext/requester"

// Check if there is a next page.
func (p *Paginator[T]) HasNext() bool {
	return p.Next != ""
}

// Check if there is a previous page.
func (p *Paginator[T]) HasPrevious() bool {
	return p.Previous != ""
}

// Default client.
func (p *Paginator[T]) client() *requester.APIClient {
	if p.Token != nil && p.Token.AccessToken != "" {
		return p.Token.Client()
	}
	return requester.NewAPIClient()
}

// Fetch the next page.
func (p *Paginator[T]) FetchNext() ([]T, error) {
	if p.Next == "" {
		return nil, nil
	}
	p.CurrentPage = p.Next
	return p.fetchResults(p.Next)
}

// Fetch the first page.
func (p *Paginator[T]) Fetch() ([]T, error) {
	p.CurrentPage = p.fetchURL
	return p.fetchResults(p.fetchURL)
}

// Fetch the previous page.
func (p *Paginator[T]) FetchPrevious() ([]T, error) {
	if p.Previous == "" {
		return nil, nil
	}
	p.CurrentPage = p.Previous
	return p.fetchResults(p.Previous)
}
