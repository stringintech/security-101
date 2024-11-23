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

func (f *JwtAuthenticationFilter) DoFilter(w http.ResponseWriter, r *http.Request, chain *FilterChain) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		chain.DoFilter(w, r)
		return
	}

	tokenString := authHeader[7:]

	username, err := f.jwtService.ValidateTokenAndGetUsername(tokenString)
	if err != nil {
		chain.DoFilter(w, r)
		return
	}

	user, exists := f.userStore.GetUserByUsername(username)
	if !exists {
		chain.DoFilter(w, r)
		return
	}

	ctx := auth.SetUserInContext(r.Context(), user)
	chain.DoFilter(w, r.WithContext(ctx))
}
