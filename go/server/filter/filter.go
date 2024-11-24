package filter

import "net/http"

// Filter defines the interface for all security filters
type Filter interface {
	// DoFilter processes the request and either:
	// 1. Calls chain.Next() to continue to the next filter
	// 2. Writes to the Response to end the chain
	DoFilter(Response http.ResponseWriter, r *http.Request, chain Chain)
}
