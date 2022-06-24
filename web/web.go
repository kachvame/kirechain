package web

import (
	"net/http"

	"github.com/fr3fou/polo/polo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
		w.Write([]byte(a.Chain.NextUntilEnd(in)))
		w.WriteHeader(200)
		return
	} else {
		w.Write([]byte(a.Chain.NextUntilEnd(a.Chain.RandomState())))
		w.WriteHeader(200)
		return
	}
}

func (a *API) healthy(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
	w.WriteHeader(200)
}

func (a *API) ready(w http.ResponseWriter, r *http.Request) {
	if a.Chain != nil {
		w.Write([]byte("ready"))
		w.WriteHeader(200)
		return
	}
	w.Write([]byte("chain not built yet"))
	w.WriteHeader(503)
}
