package filter

import (
	"net/http"
)

// FilterChain implements the actual chain of filters
type FilterChain struct {
	filters []Filter     // List of filters to process
	handler http.Handler // Final handler after all filters
}

// NewFilterChain creates a new filter chain with a final handler
func NewFilterChain(handler http.Handler, filters ...Filter) *FilterChain {
	return &FilterChain{
		filters: filters,
		handler: handler,
	}
}

// ServeHTTP implements http.Handler to make FilterChain a valid handler
func (fc *FilterChain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Start the chain processing with the first filter
	if len(fc.filters) > 0 {
		w := &ResponseWriter{
			ResponseWriter: w,
			written:        false,
		}
		chain := &RequestFilterChain{
			FilterChain: fc,
			position:    0,
		}
		fc.filters[0].DoFilter(w, r, chain)
	} else {
		// No filters, just call the final handler
		fc.handler.ServeHTTP(w, r)
	}
}
