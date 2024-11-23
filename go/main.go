package main

import (
	"github.com/stringintech/security-101/server"
	"github.com/stringintech/security-101/server/auth"
	"github.com/stringintech/security-101/server/filter"
	"github.com/stringintech/security-101/store"
	"log"
	"net/http"
	"time"
)

func main() {
	userStore := store.NewUserStore()
	jwtService := auth.NewJwtService(auth.JwtServiceConfig{
		Secret:             []byte("420a0ba4703e5392e585ec1add824e2be56ab16d89e64b7731bf870d74dd9e82"),
		ExpirationInterval: 5 * 24 * time.Hour,
	})

	userHandler := server.NewUserHandler(userStore, jwtService)

	dispatcher := server.NewDispatcher()
	dispatcher.Register("/auth/register", userHandler.Register)
	dispatcher.Register("/auth/login", userHandler.Login)
	dispatcher.Register("/users/me", userHandler.Me)

	filterChain := filter.NewFilterChain(
		dispatcher,
		filter.NewJwtAuthenticationFilter(jwtService, userStore),
		filter.NewAuthenticationFilter([]string{"/auth/"}),
	)

	log.Fatal(http.ListenAndServe(":8080", filterChain))
}
