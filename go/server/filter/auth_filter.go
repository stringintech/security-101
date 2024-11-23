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

func (f *AuthenticationFilter) DoFilter(w http.ResponseWriter, r *http.Request, chain *FilterChain) {
	for _, path := range f.publicPaths {
		if strings.HasPrefix(r.URL.Path, path) {
			chain.DoFilter(w, r)
			return
		}
	}

	_, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chain.DoFilter(w, r)
}
