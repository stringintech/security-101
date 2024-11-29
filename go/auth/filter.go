package auth

// Filter defines the interface for all security filters
type Filter interface {
	// DoFilter processes the request and either:
	// 1. Calls request.Proceed() to continue to the next filter
	// 2. Writes to the Response to end the chain TODO
	DoFilter(request *Request)
}
