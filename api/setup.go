package api

import (
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/islazy/tui"
	"github.com/evilsocket/joe/models"
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
			// GET|POST /api/v1/auth
			r.Get("/auth", api.Authenticate)
			r.Post("/auth", api.Authenticate)

			// GET /api/v1/queries/
			r.Get("/queries", api.ListQueries)

			// GET /api/v1/query/<name>/view
			r.Get("/query/{name:.+}/view", api.ShowQuery)

			// GET /api/v1/query/<name>
			r.Get("/query/{name:.+}", api.RunQuery)
			// POST /api/v1/query/<name>
			r.Post("/query/{name:.+}", api.RunQuery)

			// POST /api/v1/query/<name>/explain
			r.Post("/query/{name:.+}/explain", api.ExplainQuery)
			// GET /api/v1/query/<name>/explain
			r.Get("/query/{name:.+}/explain", api.ExplainQuery)

			// GET /api/v1/query/<name>/<view_name>
			r.Get("/query/{name:.+}/{view_name:.+}", api.RunView)
			// POST /api/v1/query/<name>/<view_name>
			r.Post("/query/{name:.+}/{view_name:.+}", api.RunView)
		})
	})

	return
}

func (api *API) Run(addr string) {
	log.Info("joe api v%s starting on %s ...", Version, addr)

	models.Queries.Range(func(key, value interface{}) bool {
		q := value.(*models.Query)
		log.Debug("  %s", tui.Dim(q.Expression))
		log.Debug("    http://%s/api/v1/query/%s(.json|csv)(/explain?)", addr, key)
		for name, _ := range q.Views {
			log.Debug("      http://%s/api/v1/query/%s/%s(.png|svg)", addr, key, name)
		}
		return true
	})

	log.Fatal("%v", http.ListenAndServe(addr, api.Router))
}
