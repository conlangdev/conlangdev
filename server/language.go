package server

import (
	"encoding/json"
	"net/http"

	"github.com/conlangdev/conlangdev"
)

func (s *Server) registerLanguageRoutes() {
	s.authenticatedRouter.HandleFunc("/language", s.handleIndexLanguage).Methods("GET")
	s.authenticatedRouter.HandleFunc("/language", s.handleCreateLanguage).Methods("POST")
}

func (s *Server) handleIndexLanguage(w http.ResponseWriter, r *http.Request) {
	// Find languages
	user := conlangdev.GetUserFromContext(r.Context())
	languages, err := s.LanguageService.FindLanguagesForUser(r.Context(), user)
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}
	// Marshal into JSON response
	response, err := json.Marshal(map[string][]*conlangdev.Language{
		"languages": languages,
	})
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}
	w.Write(response)
}

func (s *Server) handleCreateLanguage(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var create conlangdev.LanguageCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		handleError(&conlangdev.Error{
			Code:       conlangdev.EBADREQUEST,
			Message:    "malformed request body",
			StatusCode: http.StatusBadRequest,
		}).ServeHTTP(w, r)
		return
	}
	// Get user & create language
	user := conlangdev.GetUserFromContext(r.Context())
	language, err := s.LanguageService.CreateLanguageForUser(r.Context(), user, create)
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}
	// Marshal new language into JSON response
	response, err := json.Marshal(map[string]*conlangdev.Language{
		"language": language,
	})
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}
	w.Write(response)
}
