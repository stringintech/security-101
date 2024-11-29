package auth

import (
	"strings"
)

type JwtAuthenticationFilter struct {
	jwtService *JwtService
	userStore  UserStore
}

func NewJwtAuthenticationFilter(jwt *JwtService, store UserStore) *JwtAuthenticationFilter {
	return &JwtAuthenticationFilter{jwt, store}
}

func (f *JwtAuthenticationFilter) DoFilter(r *Request) {
	authHeader := r.Http.Header.Get("Authorization")

	// If no auth header or not Bearer token, continue chain
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		r.Proceed()
		return
	}

	tokenString := authHeader[7:]

	// Validate token
	username, err := f.jwtService.ValidateTokenAndGetUsername(tokenString)
	if err != nil {
		//TODO log?
		r.Proceed()
		return
	}

	// Get user
	user, exists := f.userStore.GetUserByUsername(username)
	if !exists {
		//TODO log?
		r.Proceed()
		return
	}

	// Set user in context and continue chain
	ctx := SetUserInContext(r.Http.Context(), user)
	r.Http = r.Http.WithContext(ctx) //TODO :)
	r.Proceed()
}
