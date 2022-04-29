package server

import (
	"net/http"

	"github.com/conlangdev/conlangdev"
	"github.com/gorilla/mux"
)

type Router struct {
	router *mux.Router
}

type RouteBuilder struct {
	router                *Router
	handlerFunc           func(http.ResponseWriter, *http.Request)
	authorizedHandlerFunc func(http.ResponseWriter, *http.Request, *conlangdev.User)
	requireGuest          bool
}

func NewRouter() *Router {
	return &Router{mux.NewRouter().StrictSlash(true)}
}

func (r *Router) Use(mwf ...mux.MiddlewareFunc) {
	r.router.Use(mwf...)
}

func (r *Router) Prefix(prefix string, routes func(*Router)) {
	subrouter := &Router{r.router.PathPrefix(prefix).Subrouter()}
	routes(subrouter)
}

func (r *Router) GetHandler() http.HandlerFunc {
	return http.HandlerFunc(r.router.ServeHTTP)
}

func (r *Router) HandleNotFound(f func(http.ResponseWriter, *http.Request)) {
	r.router.NotFoundHandler = http.HandlerFunc(f)
}

func (r *Router) Handle(f func(http.ResponseWriter, *http.Request)) *RouteBuilder {
	return &RouteBuilder{r, f, nil, false}
}

func (r *Router) Guest(f func(http.ResponseWriter, *http.Request)) *RouteBuilder {
	return &RouteBuilder{r, f, nil, true}
}

func (r *Router) Authorized(f func(http.ResponseWriter, *http.Request, *conlangdev.User)) *RouteBuilder {
	return &RouteBuilder{r, nil, f, false}
}

func (r *RouteBuilder) buildAuthorized(path string) *mux.Route {
	return r.router.router.HandleFunc(
		path,
		func(w http.ResponseWriter, rq *http.Request) {
			user := conlangdev.GetUserFromContext(rq.Context())
			if user == nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			r.authorizedHandlerFunc(w, rq, user)
		},
	)
}

func (r *RouteBuilder) buildGuest(path string) *mux.Route {
	return r.router.router.HandleFunc(
		path,
		func(w http.ResponseWriter, rq *http.Request) {
			user := conlangdev.GetUserFromContext(rq.Context())
			if user != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			r.handlerFunc(w, rq)
		},
	)
}

func (r *RouteBuilder) buildGeneric(path string) *mux.Route {
	return r.router.router.HandleFunc(path, r.handlerFunc)
}

func (r *RouteBuilder) build(path string) *mux.Route {
	if r.authorizedHandlerFunc != nil {
		return r.buildAuthorized(path)
	} else if r.requireGuest {
		return r.buildGuest(path)
	}
	return r.buildGeneric(path)
}

func (r *RouteBuilder) GET(path string) *mux.Route {
	return r.build(path).Methods("GET")
}

func (r *RouteBuilder) POST(path string) *mux.Route {
	return r.build(path).Methods("POST")
}
