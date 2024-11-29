package auth

import "net/http"

type Request struct {
	Http           *http.Request
	ResponseWriter http.ResponseWriter
	chain          *FilterChain // Can be shared between multiple instances since it is immutable
	position       int          // Current position in filter chain
}

func NewRequest(fc *FilterChain, w http.ResponseWriter, r *http.Request) *Request {
	w = &ResponseWriter{
		ResponseWriter: w,
		written:        false,
	}
	return &Request{
		Http:           r,
		ResponseWriter: w,
		chain:          fc,
		position:       0,
	}
}

func (r *Request) Proceed() {
	if r.ResponseWriter.(*ResponseWriter).Written() {
		return
	}

	// Move to next filter
	r.position++

	// If we have more filters, continue the chain
	if r.position < len(r.chain.filters) {
		r.chain.filters[r.position].DoFilter(r)
	}
}

// ResponseWriter wraps http.ResponseWriter to track if headers were written
type ResponseWriter struct {
	http.ResponseWriter
	written bool
}

// WriteHeader overrides http.ResponseWriter
func (w *ResponseWriter) WriteHeader(code int) {
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

// Write overrides http.ResponseWriter
func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(b)
}

// Written returns true if the response has been written to
func (w *ResponseWriter) Written() bool {
	return w.written
}
