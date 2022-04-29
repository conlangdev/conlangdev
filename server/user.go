package server

import (
	"encoding/json"
	"net/http"

	"github.com/conlangdev/conlangdev"
	"github.com/gorilla/mux"
)

func (s *Server) registerUserRoutes() {
	s.router.HandleFunc("/user/{username}", s.handleViewUser).Methods("GET")
	s.router.HandleFunc("/auth", s.handleCheckAuth).Methods("GET")
	s.unauthenticatedRouter.HandleFunc("/auth/register", s.handleRegisterUser).Methods("POST")
	s.unauthenticatedRouter.HandleFunc("/auth/login", s.handleLoginUser).Methods("POST")
}

func (s *Server) handleCheckAuth(w http.ResponseWriter, r *http.Request) {
	user := conlangdev.GetUserFromContext(r.Context())
	if user != nil {
		response, err := json.Marshal(map[string]interface{}{
			"authentication": map[string]interface{}{
				"authenticated": true,
				"user":          user,
			},
		})
		if err != nil {
			handleError(err).ServeHTTP(w, r)
			return
		}

		w.Write(response)
		return
	}

	response, err := json.Marshal(map[string]interface{}{
		"authentication": map[string]interface{}{
			"authenticated": false,
			"user":          nil,
		},
	})
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	w.Write(response)
}

func (s *Server) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var create conlangdev.UserCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		handleError(&conlangdev.Error{
			Code:       conlangdev.EBADREQUEST,
			Message:    "malformed request body",
			StatusCode: http.StatusBadRequest,
		}).ServeHTTP(w, r)
		return
	}

	user, err := s.UserService.CreateUser(r.Context(), create)
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	response, err := json.Marshal(map[string]*conlangdev.User{
		"user": user,
	})
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	w.Write(response)
}

func (s *Server) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	var loginPayload struct {
		Username string
		Password string
	}
	if err := json.NewDecoder(r.Body).Decode(&loginPayload); err != nil {
		handleError(&conlangdev.Error{
			Code:    conlangdev.EBADREQUEST,
			Message: "malformed request body",
		}).ServeHTTP(w, r)
		return
	}

	user, err := s.UserService.GetUserByUsername(r.Context(), loginPayload.Username)
	if err != nil {
		handleError(&conlangdev.Error{
			Code:       conlangdev.EUNAUTHORIZED,
			Message:    "incorrect username or password",
			StatusCode: http.StatusUnauthorized,
		}).ServeHTTP(w, r)
		return
	}

	if err := s.UserService.CheckUserPassword(r.Context(), user, loginPayload.Password); err != nil {
		handleError(&conlangdev.Error{
			Code:       conlangdev.EUNAUTHORIZED,
			Message:    "incorrect username or password",
			StatusCode: http.StatusUnauthorized,
		}).ServeHTTP(w, r)
		return
	}

	jwt, err := s.UserService.GenerateJWTForUser(r.Context(), user)
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	response, err := json.Marshal(map[string]interface{}{
		"authentication": map[string]interface{}{
			"jwt":  jwt,
			"user": user,
		},
	})
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	w.Write(response)
}

func (s *Server) handleViewUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	user, err := s.UserService.GetUserByUsername(r.Context(), username)
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	userView, err := s.UserService.GetViewForUser(r.Context(), user)
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	response, err := json.Marshal(userView)
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	w.Write(response)
}
