package server

import (
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Dispatcher struct {
	// Maps URL patterns to their handlers
	handlers map[string]HandlerFunc
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string]HandlerFunc),
	}
}

// Register adds a new handler for a URL pattern
func (d *Dispatcher) Register(pattern string, handler HandlerFunc) {
	d.handlers[pattern] = handler
}

// ServeHTTP routes the request to the appropriate handler
func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Find exact match
	if handler, exists := d.handlers[r.URL.Path]; exists {
		handler(w, r)
		return
	}
	http.NotFound(w, r)
}
