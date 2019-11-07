package api

import (
	"github.com/evilsocket/joe/models"
	"github.com/go-chi/chi"
	"github.com/wcharczuk/go-chart"
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

// GET /api/v1/query/<name>/view
func (api *API) ShowQuery(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if q, found := models.Queries.Load(name); found {
		JSON(w, http.StatusOK, q)
	} else {
		ERROR(w, http.StatusNotFound, ErrEmpty)
	}
}

// GET|POST /api/v1/query/<name>
func (api *API) RunQuery(w http.ResponseWriter, r *http.Request) {
	name, ext := parseName("name", r)
	if q, found := models.Queries.Load(name); found {
		params := parseParameters(r)
		if rows, err := q.(*models.Query).Query(params); err != nil {
			ERROR(w, http.StatusBadRequest, err)
		} else if ext == "csv" {
			CSV(w, http.StatusOK, rows)
		} else if ext == "json" {
			JSON(w, http.StatusOK, rows)
		} else {
			ERROR(w, http.StatusNotFound, ErrEmpty)
		}
	} else {
		ERROR(w, http.StatusNotFound, ErrEmpty)
	}
}

// GET|POST /api/v1/query/<name>/explain
func (api *API) ExplainQuery(w http.ResponseWriter, r *http.Request) {
	name, _ := parseName("name", r)
	if q, found := models.Queries.Load(name); found {
		params := parseParameters(r)
		if rows, err := q.(*models.Query).Explain(params); err != nil {
			ERROR(w, http.StatusBadRequest, err)
		} else {
			JSON(w, http.StatusOK, rows)
		}
	} else {
		ERROR(w, http.StatusNotFound, ErrEmpty)
	}
}

// GET|POST /api/v1/query/<name>/<view_name>
func (api *API) RunView(w http.ResponseWriter, r *http.Request) {
	name, _ := parseName("name", r)
	if q, found := models.Queries.Load(name); found {
		query := q.(*models.Query)
		viewName, viewExt := parseName("view_name", r)
		view := query.View(viewName)
		if view == nil {
			ERROR(w, http.StatusNotFound, ErrEmpty)
			return
		}

		params := parseParameters(r)
		rows, err := q.(*models.Query).Query(params)
		if err != nil {
			ERROR(w, http.StatusBadRequest, err)
			return
		}

		graph := view.Call(rows)

		if viewExt == "png" {
			w.Header().Set("Content-Type", "image/png")
			w.WriteHeader(http.StatusOK)
			graph.Render(chart.PNG, w)
		} else if viewExt == "svg" {
			w.Header().Set("Content-Type", "image/svg")
			w.WriteHeader(http.StatusOK)
			graph.Render(chart.SVG, w)
		} else {
			ERROR(w, http.StatusNotFound, ErrEmpty)
		}
	} else {
		ERROR(w, http.StatusNotFound, ErrEmpty)
	}
}