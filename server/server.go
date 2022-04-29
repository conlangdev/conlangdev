package server

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/conlangdev/conlangdev"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	server   *http.Server
	listener net.Listener
	router   *Router

	Addr string

	UserService     conlangdev.UserService
	LanguageService conlangdev.LanguageService
	WordService     conlangdev.WordService
}

func NewServer() *Server {
	// Set up server object with base routers
	server := &Server{
		server: &http.Server{},
		router: NewRouter(),
	}

	// Add middleware
	server.router.Use(contentTypeJSON)
	server.router.Use(logger)
	server.router.Use(server.authenticate)

	// Register routes
	server.registerUserRoutes()
	server.registerLanguageRoutes()
	server.registerWordRoutes()

	// Allocate handler to our router and return server
	server.server.Handler = server.router.GetHandler()
	server.router.HandleNotFound(handleNotFound)
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
		w.WriteHeader(http.StatusNotFound)
		if response, err := json.Marshal(map[string]interface{}{
			"error": map[string]string{
				"code":    "not_found",
				"message": "could not find the given resource",
			},
		}); err == nil {
			w.Write(response)
		}
	})).ServeHTTP(w, r)
}

// Handles an error from the business layer, writing it to the HTTP response
// in a client-friendly way.
//
// Errors should be given as `*conlangdev.Error` or `*conlangdev.FieldsError`
// in order to be handled properly, otherwise a generic 500 Internal Server
// Error will be returned.
func handleError(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a conlangdev error repsonse.
		if cd_err, ok := err.(*conlangdev.Error); ok {
			w.WriteHeader(cd_err.StatusCode)
			if response, err := json.Marshal(map[string]interface{}{
				"error": map[string]string{
					"code":    cd_err.Code,
					"message": cd_err.Message,
				},
			}); err == nil {
				w.Write(response)
			}
			return
		}

		// Return a conlangdev fields error response.
		if cd_fields, ok := err.(*conlangdev.FieldsError); ok {
			w.WriteHeader(cd_fields.StatusCode)
			if response, err := json.Marshal(map[string]interface{}{
				"error": map[string]interface{}{
					"code":    cd_fields.Code,
					"message": cd_fields.Message,
					"fields":  cd_fields.Fields,
				},
			}); err == nil {
				w.Write(response)
			}
			return
		}

		// Unknown error (not passed along as conlangdev.Error or
		// conlangdev.FieldsError). Log it & return a generic 500
		// Internal Server Error response.
		log.Errorf("unknown error occurred whilst handling request: (type %T) %v", err, err)
		w.WriteHeader(http.StatusInternalServerError)
		// Marshal and write error response as JSON
		if response, err := json.Marshal(map[string]interface{}{
			"error": map[string]string{
				"code":    "server_error",
				"message": "something unexpected went wrong",
			},
		}); err == nil {
			w.Write(response)
		}
	})
}

// Middleware function which validates any user authentication in the
// request and places the authenticated user into the context, if any.
//
// This middleware does *not* specify whether a user must or must not
// be authenticated to continue. For that, see the route-building
// functions `Router.Authorized()` and `Router.Guest()`
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

func (s *Server) WithWordService(ws conlangdev.WordService) *Server {
	s.WordService = ws
	return s
}
