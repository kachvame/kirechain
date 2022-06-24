package web

import (
	"net/http"

	"github.com/fr3fou/polo/polo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type API struct {
	*chi.Mux
	Chain *polo.Chain
}

func New(c *polo.Chain) *API {
	a := &API{
		Mux:   chi.NewRouter(),
		Chain: c,
	}

	a.Use(middleware.RequestID)
	a.Use(middleware.RealIP)
	a.Use(middleware.Logger)
	a.Use(middleware.Recoverer)

	a.Get("/", a.gen)
	a.Get("/healthy", a.healthy)
	a.Get("/ready", a.ready)

	return a
}

func (a *API) gen(w http.ResponseWriter, r *http.Request) {
	in := r.URL.Query().Get("in")
	if in != "" {
		render.PlainText(w, r, a.Chain.NextUntilEnd(in))
	} else {
		render.PlainText(w, r, a.Chain.NextUntilEnd(a.Chain.RandomState()))
	}
}

func (a *API) healthy(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "ok")
}

func (a *API) ready(w http.ResponseWriter, r *http.Request) {
	if a.Chain != nil {
		render.PlainText(w, r, "chain is ready")
		return
	}
	render.Status(r, 503)
	render.PlainText(w, r, "chain is not ready")
}
