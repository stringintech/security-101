package server

import (
	"github.com/stringintech/security-101/auth"
	"github.com/stringintech/security-101/store"
	"log"
	"net/http"
	"runtime/debug"
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
	w = &ResponseWriter{ResponseWriter: w}

	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			log.Printf("panic: %v\n%s", err, stack)

			if rw, ok := w.(*ResponseWriter); !ok || !rw.Written() { //TODO? we should be able to override any written stuff instead
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}
	}()

	s.DoFilter(w, r)
}

func (s *Server) Serve() {
	log.Fatal(http.ListenAndServe(":8080", s)) //TODO make configurable
}
