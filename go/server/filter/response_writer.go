package filter

import "net/http"

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
