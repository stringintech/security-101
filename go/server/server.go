package server

import (
	"github.com/stringintech/security-101/server/auth"
	"github.com/stringintech/security-101/server/filter"
	"github.com/stringintech/security-101/store"
	"net/http"
)

func New(userStore *store.UserStore, jwtService *auth.JwtService) http.Handler {
	userHandler := NewUserHandler(userStore, jwtService)

	dispatcher := NewDispatcher()
	dispatcher.Register("/auth/register", userHandler.Register)
	dispatcher.Register("/auth/login", userHandler.Login)
	dispatcher.Register("/users/me", userHandler.Me)

	filterChain := filter.NewFilterChain(
		dispatcher,
		filter.NewJwtAuthenticationFilter(jwtService, userStore),
		filter.NewAuthenticationFilter([]string{"/auth/"}),
	)

	return filterChain
}
