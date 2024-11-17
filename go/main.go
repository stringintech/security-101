package main

import (
	"github.com/stringintech/security-101/server"
	"github.com/stringintech/security-101/server/auth"
	"log"
	"net/http"
	"time"

	"github.com/stringintech/security-101/store"
)

func main() {
	// Initialize dependencies
	userStore := store.NewUserStore()
	jwtService := auth.NewJwtService(auth.JwtServiceConfig{
		Secret:             []byte("420a0ba4703e5392e585ec1add824e2be56ab16d89e64b7731bf870d74dd9e82"),
		ExpirationInterval: 5 * 24 * time.Hour,
	})

	// Initialize handlers
	userHandler := server.NewUserHandler(userStore, jwtService)

	// Public endpoints
	http.HandleFunc("/auth/register", userHandler.Register)
	http.HandleFunc("/auth/login", userHandler.Login)

	// Protected endpoint
	authMiddleware := auth.NewMiddleware(userStore, jwtService)
	http.HandleFunc("/users/me", authMiddleware.WrapHandler(userHandler.Me))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
