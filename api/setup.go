package api

import (
	"github.com/evilsocket/islazy/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

type API struct {
	Router *chi.Mux
}

func Setup() (err error, api *API) {
	api = &API{
		Router: chi.NewRouter(),
	}

	api.Router.Use(CORS)

	api.Router.Use(middleware.DefaultCompress)

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// GET /api/v1/queries/
			r.Get("/queries", api.ListQueries)
			// GET /api/v1/query/<name>
			r.Get("/query/{name:.+}", api.ShowQuery)
			// POST /api/v1/query/<name>
			r.Post("/query/{name:.+}", api.RunQuery)
			// POST /api/v1/query/<name>/explain
			r.Post("/query/{name:.+}/explain", api.ExplainQuery)
		})
	})

	return
}

func (api *API) Run(addr string) {
	log.Info("joe api starting on %s ...", addr)
	log.Fatal("%v", http.ListenAndServe(addr, api.Router))
}
