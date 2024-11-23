package filter

import (
	"net/http"
)

// Filter defines the interface for all security filters
type Filter interface {
	// DoFilter processes the request and either:
	// 1. Calls chain.Next() to continue to the next filter
	// 2. Writes to the response to end the chain
	DoFilter(w http.ResponseWriter, r *http.Request, chain Chain)
}

// Chain defines how filters are chained together
type Chain interface {
	// Next continues processing with the next filter or final handler
	Next(w http.ResponseWriter, r *http.Request)
}

// FilterChain implements the actual chain of filters
type FilterChain struct {
	filters  []Filter     // List of filters to process
	handler  http.Handler // Final handler after all filters
	position int          // Current position in filter chain
}

// NewFilterChain creates a new filter chain with a final handler
func NewFilterChain(handler http.Handler, filters ...Filter) *FilterChain {
	return &FilterChain{
		filters:  filters,
		handler:  handler,
		position: 0,
	}
}

// ServeHTTP implements http.Handler to make FilterChain a valid handler
func (fc *FilterChain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Start the chain processing with the first filter
	if len(fc.filters) > 0 {
		fc.filters[0].DoFilter(w, r, fc)
	} else {
		// No filters, just call the final handler
		fc.handler.ServeHTTP(w, r)
	}
}

// Next implements Chain interface
func (fc *FilterChain) Next(w http.ResponseWriter, r *http.Request) {
	rw, ok := w.(*ResponseWriter)
	if !ok {
		// Wrap the original ResponseWriter if not already wrapped
		rw = &ResponseWriter{ResponseWriter: w}
		w = rw
	} else if rw.Written() {
		return
	}

	// Move to next filter
	fc.position++

	// If we have more filters, continue the chain
	if fc.position < len(fc.filters) {
		fc.filters[fc.position].DoFilter(w, r, fc)
		return
	}

	// No more filters, call the final handler
	fc.handler.ServeHTTP(w, r)
}
