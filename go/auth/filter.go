package auth

import "net/http"

type Filter interface {
	DoFilter(w http.ResponseWriter, r *http.Request, filterChain FilterChain)
}
