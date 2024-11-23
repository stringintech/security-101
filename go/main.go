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

	filterChain := filter.NewFilterChain(
		filter.NewJwtAuthenticationFilter(jwtService, userStore),
		filter.NewAuthenticationFilter([]string{"/auth/"}),
	)

	// Map all requests through the filter chain
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Special case direct handling for register/login endpoints
		switch r.URL.Path {
		case "/auth/register":
			//TODO
			//unprotected endpoints should also go through the filters;
			//also consider adding a dispatcher component
			userHandler.Register(w, r)
			return
		case "/auth/login":
			userHandler.Login(w, r)
			return
		case "/users/me":
			filterChain.DoFilter(w, r)
			if r.Method == http.MethodPost { //FIXME if authentication fails we won't proceed
				userHandler.Me(w, r)
			}
			return
		default:
			http.NotFound(w, r)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
