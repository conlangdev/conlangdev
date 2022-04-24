package server

import "net/http"

func (s *Server) registerLanguageRoutes() {
	s.authenticatedRouter.HandleFunc("/language", s.handleLanguageIndex).Methods("GET")
}

func (s *Server) handleLanguageIndex(w http.ResponseWriter, r *http.Request) {
}
