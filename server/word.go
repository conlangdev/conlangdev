package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/conlangdev/conlangdev"
	"github.com/gorilla/mux"
)

func (s *Server) registerWordRoutes() {
	s.router.Prefix("/word/{username}/{language}", func(word *Router) {
		word.Handle(s.handleIndexWord).GET("")
		word.Authorized(s.handleCreateWord).GET("")
		word.Handle(s.handleViewWord).GET("/{word}")
	})
}

func (s *Server) handleIndexWord(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userx, err := s.UserService.GetUserByUsername(r.Context(), params["username"])
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	language, err := s.LanguageService.GetLanguageByUserAndSlug(r.Context(), userx, params["language"])
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	words, err := s.WordService.FindWordsForLanguage(r.Context(), language)
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	response, err := json.Marshal(map[string][]*conlangdev.Word{
		"words": words,
	})
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	w.Write(response)
}

func (s *Server) handleCreateWord(w http.ResponseWriter, r *http.Request, user *conlangdev.User) {
	params := mux.Vars(r)
	userx, err := s.UserService.GetUserByUsername(r.Context(), params["username"])
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}
	if user.ID != userx.ID {
		handleError(&conlangdev.Error{
			Code:       conlangdev.EUNAUTHORIZED,
			Message:    "you must be the owner of a language to add words",
			StatusCode: http.StatusForbidden,
		}).ServeHTTP(w, r)
		return
	}

	language, err := s.LanguageService.GetLanguageByUserAndSlug(r.Context(), userx, params["language"])
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	var create conlangdev.WordCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		handleError(&conlangdev.Error{
			Code:       conlangdev.EBADREQUEST,
			Message:    "malformed request body",
			StatusCode: http.StatusBadRequest,
		}).ServeHTTP(w, r)
		return
	}

	word, err := s.WordService.CreateWordForLanguage(r.Context(), language, create)
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	response, err := json.Marshal(map[string]*conlangdev.Word{
		"word": word,
	})
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	w.Write(response)
}

func (s *Server) handleViewWord(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userx, err := s.UserService.GetUserByUsername(r.Context(), params["username"])
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	language, err := s.LanguageService.GetLanguageByUserAndSlug(r.Context(), userx, params["language"])
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	wordUID, err := strconv.ParseUint(params["word"], 10, 64)
	if err != nil {
		handleError(&conlangdev.Error{
			Code:       conlangdev.EBADREQUEST,
			Message:    "invalid word ID",
			StatusCode: http.StatusBadRequest,
		}).ServeHTTP(w, r)
		return
	}

	word, err := s.WordService.GetWordByLanguageAndUID(r.Context(), language, uint(wordUID))
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	response, err := json.Marshal(map[string]*conlangdev.Word{
		"word": word,
	})
	if err != nil {
		handleError(err).ServeHTTP(w, r)
		return
	}

	w.Write(response)
}
