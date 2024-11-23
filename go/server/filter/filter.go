package filter

import (
	"net/http"
)

type FilterChain struct {
	filters []Filter
}

type Filter interface {
	DoFilter(w http.ResponseWriter, r *http.Request, chain *FilterChain)
}

func NewFilterChain(filters ...Filter) *FilterChain {
	return &FilterChain{filters: filters}
}

func (c *FilterChain) DoFilter(w http.ResponseWriter, r *http.Request) {
	if len(c.filters) > 0 {
		filter := c.filters[0]
		remainingFilters := append([]Filter{}, c.filters[1:]...)
		remainingChain := &FilterChain{filters: remainingFilters}
		filter.DoFilter(w, r, remainingChain)
	}
}
