package auth

import (
	"context"
	"net/http"
)

type contextKey string

const userContextKey contextKey = "user"

type Middleware struct {
	store UserStore
	jwt   *JwtService
}

func NewMiddleware(store UserStore, jwt *JwtService) *Middleware {
	return &Middleware{
		store: store,
		jwt:   jwt,
	}
}

func (m *Middleware) WrapHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[7:]
		username, err := m.jwt.ValidateTokenAndGetUsername(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, exists := m.store.GetUserByUsername(username)
		if !exists {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetUserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userContextKey).(User)
	return user, ok
}
