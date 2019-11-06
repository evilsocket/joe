package api

import (
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/joe/models"
	"github.com/go-chi/chi"
	"net/http"
)

// GET /api/v1/queries/
func (api *API) ListQueries(w http.ResponseWriter, r *http.Request) {
	queries := make([]*models.Query, 0)

	models.Queries.Range(func(key, value interface{}) bool {
		queries = append(queries, value.(*models.Query))
		return true
	})

	JSON(w, http.StatusOK, queries)
}

// GET /api/v1/query/<name>
func (api *API) ShowQuery(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if q, found := models.Queries.Load(name); found {
		JSON(w, http.StatusOK, q)
	} else {
		ERROR(w, http.StatusNotFound, ErrEmpty)
	}
}

// POST /api/v1/query/<name>
func (api *API) RunQuery(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if q, found := models.Queries.Load(name); found {
		params := make(map[string]interface{})
		if err := r.ParseForm(); err == nil {
			for key, values := range r.PostForm {
				params[key] = values[0]
			}
		} else {
			log.Warning("%v", err)
		}

		if rows, err := q.(*models.Query).Query(params); err != nil {
			ERROR(w, http.StatusBadRequest, err)
		} else {
			JSON(w, http.StatusOK, rows)
		}
	} else {
		ERROR(w, http.StatusNotFound, ErrEmpty)
	}
}

// POST /api/v1/query/<name>/explain
func (api *API) ExplainQuery(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if q, found := models.Queries.Load(name); found {
		params := make(map[string]interface{})
		if err := r.ParseForm(); err == nil {
			for key, values := range r.PostForm {
				params[key] = values[0]
			}
		} else {
			log.Warning("%v", err)
		}

		if rows, err := q.(*models.Query).Explain(params); err != nil {
			ERROR(w, http.StatusBadRequest, err)
		} else {
			JSON(w, http.StatusOK, rows)
		}
	} else {
		ERROR(w, http.StatusNotFound, ErrEmpty)
	}
}