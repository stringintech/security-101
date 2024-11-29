package auth

import (
	"net/http"
	"strings"
)

type AuthenticationFilter struct {
	publicPaths []string
}

func NewAuthenticationFilter(publicPaths []string) *AuthenticationFilter {
	return &AuthenticationFilter{publicPaths}
}

func (f *AuthenticationFilter) DoFilter(r *Request) {
	// Continue chain if path is public
	for _, path := range f.publicPaths {
		if strings.HasPrefix(r.Http.URL.Path, path) {
			r.Proceed()
			return
		}
	}

	// Check if user is authenticated
	_, ok := GetUserFromContext(r.Http.Context())
	if !ok {
		http.Error(r.ResponseWriter, "Unauthorized", http.StatusUnauthorized) //TODO add method to r?
		return
	}

	// User is authenticated, continue chain
	r.Proceed()
}
