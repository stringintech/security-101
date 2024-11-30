package auth

import (
	"net/http"
	"strings"
)

type JwtFilter struct {
	jwtService *JwtService
	userStore  UserStore
}

func NewJwtFilter(jwt *JwtService, store UserStore) *JwtFilter {
	return &JwtFilter{jwt, store}
}

func (f *JwtFilter) DoFilter(w http.ResponseWriter, r *http.Request, filterChain FilterChain) {
	authHeader := r.Header.Get("Authorization")

	// If no auth header or not Bearer token, continue chain
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		filterChain.DoFilter(w, r)
		return
	}

	tokenString := authHeader[7:]

	// Validate token
	username, err := f.jwtService.ValidateTokenAndGetUsername(tokenString)
	if err != nil {
		//TODO log
		filterChain.DoFilter(w, r)
		return
	}

	// Get user
	user, exists := f.userStore.GetUserByUsername(username)
	if !exists {
		//TODO log
		filterChain.DoFilter(w, r)
		return
	}

	// Set user in context and continue chain
	ctx := SetUserInContext(r.Context(), user)
	filterChain.DoFilter(w, r.WithContext(ctx))
}
