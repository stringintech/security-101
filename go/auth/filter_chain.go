package auth

import (
	"net/http"
)

type FilterChain struct {
	filters []Filter
}

func NewFilterChain(filters ...Filter) *FilterChain {
	return &FilterChain{
		filters: filters,
	}
}

func (fc *FilterChain) Filter(w http.ResponseWriter, r *http.Request) *Request {
	if len(fc.filters) == 0 {
		return nil
	}
	// Start the chain processing with the first filter
	request := NewRequest(fc, w, r)
	fc.filters[0].DoFilter(request)
	return request
}
