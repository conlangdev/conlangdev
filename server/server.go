package server

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/conlangdev/conlangdev"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	server   *http.Server
	listener net.Listener

	router                *mux.Router
	authenticatedRouter   *mux.Router
	unauthenticatedRouter *mux.Router

	Addr string

	UserService     conlangdev.UserService
	LanguageService conlangdev.LanguageService
}

func NewServer() *Server {
	// Set up server object with base routers
	router := mux.NewRouter().StrictSlash(true)
	server := &Server{
		server:                &http.Server{},
		router:                router,
		authenticatedRouter:   router.PathPrefix("/").Subrouter(),
		unauthenticatedRouter: router.PathPrefix("/").Subrouter(),
	}

	// Add middleware
	server.router.Use(contentTypeJSON)
	server.router.Use(logger)
	server.authenticatedRouter.Use(server.requireAuthentication)
	server.unauthenticatedRouter.Use(server.requireNoAuthentication)

	// Register routes
	server.registerUserRoutes()
	server.registerLanguageRoutes()

	// Allocate handler to our router and return server
	server.server.Handler = http.HandlerFunc(server.router.ServeHTTP)
	server.router.NotFoundHandler = http.HandlerFunc(handleNotFound)
	return server
}

func contentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})).ServeHTTP(w, r)
}

func (s *Server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
			jwt := strings.TrimPrefix(h, "Bearer ")
			user, err := s.UserService.GetUserByJWT(r.Context(), jwt)
			if err == nil && user != nil {
				r = r.WithContext(conlangdev.NewContextWithUser(r.Context(), user))
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := conlangdev.GetUserFromContext(r.Context()); user != nil {
			next.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	})
}

func (s *Server) requireNoAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func (s *Server) Open() (err error) {
	if s.Addr == "" {
		log.Println("No address was set, defaulting to :8000")
		s.Addr = ":8000"
	}
	if s.listener, err = net.Listen("tcp", s.Addr); err != nil {
		return err
	}
	go s.server.Serve(s.listener)
	return nil
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) WithAddr(addr string) *Server {
	s.Addr = addr
	return s
}

func (s *Server) WithUserService(us conlangdev.UserService) *Server {
	s.UserService = us
	return s
}

func (s *Server) WithLanguageService(ls conlangdev.LanguageService) *Server {
	s.LanguageService = ls
	return s
}
