package auth

import (
	"context"
	"net/http"

	"github.com/stringintech/security-101/model"
	"github.com/stringintech/security-101/store"
)

type contextKey string

const userContextKey contextKey = "user"

func Middleware(next http.HandlerFunc, store *store.UserStore, auth *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[7:]
		username, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, exists := store.FindByUsername(username)
		if !exists {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func UserFromContext(ctx context.Context) (model.User, bool) {
	user, ok := ctx.Value(userContextKey).(model.User)
	return user, ok
}
