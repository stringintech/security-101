package filter

import "net/http"

type Chain interface {
	Next(w http.ResponseWriter, r *http.Request)
}

type RequestFilterChain struct {
	*FilterChain     // Can be shared between multiple instances since it is immutable
	position     int // Current position in filter chain
}

func (fc *RequestFilterChain) Next(w http.ResponseWriter, r *http.Request) {
	if w.(*ResponseWriter).Written() {
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
