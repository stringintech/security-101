package server

import (
	"github.com/stringintech/security-101/auth"
	"github.com/stringintech/security-101/store"
	"log"
	"net/http"
)

type Server struct {
	AuthFilterChain *auth.FilterChain
	Dispatcher      http.Handler
}

func New(userStore *store.UserStore, jwtService *auth.JwtService) *Server {
	userHandler := NewUserHandler(userStore, jwtService)

	dispatcher := NewDispatcher()
	dispatcher.Register("/auth/register", userHandler.Register)
	dispatcher.Register("/auth/login", userHandler.Login)
	dispatcher.Register("/users/me", userHandler.Me)

	filterChain := auth.NewFilterChain(auth.NewJwtAuthenticationFilter(jwtService, userStore), auth.NewAuthenticationFilter([]string{"/auth/"}))

	return &Server{
		AuthFilterChain: filterChain,
		Dispatcher:      dispatcher,
	}
}

// ServeHTTP implements http.Handler to make Server a valid handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	modifiedRequest := s.AuthFilterChain.Filter(w, r) //TODO well let's think
	s.Dispatcher.ServeHTTP(modifiedRequest.ResponseWriter, modifiedRequest.Http)
}

func (s *Server) Serve() {
	log.Fatal(http.ListenAndServe(":8080", s)) //TODO make configurable
}
