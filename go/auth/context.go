package auth

import "context"

type contextKey string

const userContextKey contextKey = "user"

func GetUserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userContextKey).(User)
	return user, ok
}

func SetUserInContext(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}
