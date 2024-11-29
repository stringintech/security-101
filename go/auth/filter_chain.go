package auth

import (
	"net/http"
)

type FilterChain interface {
	DoFilter(w http.ResponseWriter, r *http.Request)
}

type FilterChainImpl struct {
	filters []Filter
}

func NewFilterChain(filters ...Filter) FilterChain {
	return &FilterChainImpl{
		filters: filters,
	}
}

func (fc *FilterChainImpl) DoFilter(w http.ResponseWriter, r *http.Request) {
	NewVirtualFilterChain(fc.filters).DoFilter(w, r)
}
