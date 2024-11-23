package filter

import (
	"github.com/stringintech/security-101/server/auth"
	"net/http"
	"strings"
)

type AuthenticationFilter struct {
	publicPaths []string
}

func NewAuthenticationFilter(publicPaths []string) *AuthenticationFilter {
	return &AuthenticationFilter{publicPaths}
}

func (f *AuthenticationFilter) DoFilter(w http.ResponseWriter, r *http.Request, chain Chain) {
	// Continue chain if path is public
	for _, path := range f.publicPaths {
		if strings.HasPrefix(r.URL.Path, path) {
			chain.Next(w, r)
			return
		}
	}

	// Check if user is authenticated
	_, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// User is authenticated, continue chain
	chain.Next(w, r)
}
