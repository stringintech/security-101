package filter

import (
	"github.com/stringintech/security-101/server/auth"
	"net/http"
	"strings"
)

type JwtAuthenticationFilter struct {
	jwtService *auth.JwtService
	userStore  auth.UserStore
}

func NewJwtAuthenticationFilter(jwt *auth.JwtService, store auth.UserStore) *JwtAuthenticationFilter {
	return &JwtAuthenticationFilter{jwt, store}
}

func (f *JwtAuthenticationFilter) DoFilter(w http.ResponseWriter, r *http.Request, chain Chain) {
	authHeader := r.Header.Get("Authorization")

	// If no auth header or not Bearer token, continue chain
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		chain.Next(w, r)
		return
	}

	tokenString := authHeader[7:]

	// Validate token
	username, err := f.jwtService.ValidateTokenAndGetUsername(tokenString)
	if err != nil {
		//TODO log?
		chain.Next(w, r)
		return
	}

	// Get user
	user, exists := f.userStore.GetUserByUsername(username)
	if !exists {
		//TODO log?
		chain.Next(w, r)
		return
	}

	// Set user in context and continue chain
	ctx := auth.SetUserInContext(r.Context(), user)
	r = r.WithContext(ctx)
	chain.Next(w, r)
}
