package main

import (
	"github.com/stringintech/security-101/auth"
	"github.com/stringintech/security-101/server"
	"github.com/stringintech/security-101/store"
	"time"
)

func main() {
	userStore := store.NewUserStore()
	jwtService := auth.NewJwtService(auth.JwtServiceConfig{
		Secret:             []byte("420a0ba4703e5392e585ec1add824e2be56ab16d89e64b7731bf870d74dd9e82"),
		ExpirationInterval: 5 * 24 * time.Hour,
	})
	s := server.New(userStore, jwtService)
	s.Serve()
}
