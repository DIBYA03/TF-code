package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
)

const (
	defaultCertLocation = "./ssl/cert.pem"
	defaultKeyLocation  = "./ssl/key.pem"

	defaultHealthCheckPath = "/healthcheck.html"
)

//Router interface, a subset of chi with some convenience methods
type Router interface {
	Delete(string, http.HandlerFunc)
	Get(string, http.HandlerFunc)
	Head(string, http.HandlerFunc)
	Options(string, http.HandlerFunc)
	Patch(string, http.HandlerFunc)
	Post(string, http.HandlerFunc)
	Put(string, http.HandlerFunc)
	Route(string, func(r chi.Router)) chi.Router
	Handle(string, http.Handler)

	HandleHealthCheck()
	ListenAndServeTLS() error
}

type router struct {
	chi *chi.Mux
}

//NewRouter returns a router. It appends some middleware as well
func NewRouter() Router {
	chi := chi.NewRouter()

	m := NewMiddleWare()

	chi.Use(chiMiddleware.RequestID)

	chi.Use(m.Logger)

	return router{
		chi: chi,
	}
}

func (r router) Delete(p string, h http.HandlerFunc) {
	r.chi.Delete(p, h)
}

func (r router) Get(p string, h http.HandlerFunc) {
	r.chi.Get(p, h)
}

func (r router) Head(p string, h http.HandlerFunc) {
	r.chi.Head(p, h)
}

func (r router) Options(p string, h http.HandlerFunc) {
	r.chi.Options(p, h)
}

func (r router) Patch(p string, h http.HandlerFunc) {
	r.chi.Patch(p, h)
}

func (r router) Post(p string, h http.HandlerFunc) {
	r.chi.Post(p, h)

}

func (r router) Put(p string, h http.HandlerFunc) {
	r.chi.Put(p, h)
}

func (r router) Route(p string, fn func(r chi.Router)) chi.Router {
	return r.chi.Route(p, fn)
}

func (r router) Handle(p string, h http.Handler) {
	r.chi.Handle(p, h)
}

func (r router) HandleHealthCheck() {
	r.chi.Get(defaultHealthCheckPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})
}

func (r router) ListenAndServeTLS() error {
	cp := os.Getenv("CONTAINER_LISTEN_PORT")

	return http.ListenAndServeTLS(fmt.Sprintf(":%s", cp), defaultCertLocation, defaultKeyLocation, r.chi)
}
