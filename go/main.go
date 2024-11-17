package main

import (
	"log"
	"net/http"

	"github.com/stringintech/security-101/auth"
	"github.com/stringintech/security-101/handler"
	"github.com/stringintech/security-101/store"
)

func main() {
	// Initialize dependencies
	userStore := store.NewUserStore()
	authService := auth.NewService([]byte("420a0ba4703e5392e585ec1add824e2be56ab16d89e64b7731bf870d74dd9e82"))

	// Initialize handlers
	userHandler := handler.NewUserHandler(userStore, authService)

	// Public endpoints
	http.HandleFunc("/auth/register", userHandler.Register)
	http.HandleFunc("/auth/login", userHandler.Login)

	// Protected endpoint
	http.HandleFunc("/users/me", auth.Middleware(userHandler.Me, userStore, authService))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
