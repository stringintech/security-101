package server

import (
	"github.com/stringintech/security-101/auth"
	"github.com/stringintech/security-101/store"
	"log"
	"net/http"
)

type Server struct {
	auth.FilterChain
}

func New(userStore *store.UserStore, jwtService *auth.JwtService) *Server {
	userHandler := NewUserHandler(userStore, jwtService)

	dispatcher := NewDispatcher()
	dispatcher.Register("/auth/register", userHandler.Register)
	dispatcher.Register("/auth/login", userHandler.Login)
	dispatcher.Register("/users/me", userHandler.Me)

	filterChain := auth.NewFilterChain(
		auth.NewJwtFilter(jwtService, userStore),
		auth.NewAuthenticationFilter([]string{"/auth/"}),
		dispatcher,
	)

	return &Server{filterChain}
}

// ServeHTTP implements http.Handler to make Server a valid handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.DoFilter(w, r) //TODO well let's think
}

func (s *Server) Serve() {
	log.Fatal(http.ListenAndServe(":8080", s)) //TODO make configurable
}
