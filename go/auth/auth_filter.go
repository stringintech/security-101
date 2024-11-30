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

func (f *AuthenticationFilter) DoFilter(w http.ResponseWriter, r *http.Request, filterChain FilterChain) {
	// Continue chain if path is public
	for _, path := range f.publicPaths {
		if strings.HasPrefix(r.URL.Path, path) {
			filterChain.DoFilter(w, r)
			return
		}
	}

	// Check if user is authenticated
	_, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// User is authenticated, continue chain
	filterChain.DoFilter(w, r)
}
