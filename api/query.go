package api

import (
	"github.com/evilsocket/joe/models"
	"github.com/wcharczuk/go-chart"
	"net/http"
)

// GET /api/v1/queries/
func (api *API) ListQueries(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if user == nil {
		ERROR(w, http.StatusUnauthorized, ErrEmpty)
		return
	}

	queries := make([]*models.Query, 0)
	models.Queries.Range(func(key, value interface{}) bool {
		query := value.(*models.Query)
		if query.Authorized(user) {
			queries = append(queries, query)
		}
		return true
	})

	JSON(w, http.StatusOK, queries)
}

// GET /api/v1/query/<name>/view
func (api *API) ShowQuery(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	name, _ := parseName("name", r)
	if q := models.FindQuery(name); q == nil {
		ERROR(w, http.StatusNotFound, ErrEmpty)
	} else if q.Authorized(user) == false {
		ERROR(w, http.StatusUnauthorized, ErrEmpty)
	} else {
		JSON(w, http.StatusOK, q)
	}
}

// GET|POST /api/v1/query/<name>
func (api *API) RunQuery(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	name, ext := parseName("name", r)
	if q := models.FindQuery(name); q == nil {
		ERROR(w, http.StatusNotFound, ErrEmpty)
	} else if q.Authorized(user) == false {
		ERROR(w, http.StatusUnauthorized, ErrEmpty)
	} else {
		params := parseParameters(r)
		if rows, err := q.Query(params); err != nil {
			ERROR(w, http.StatusBadRequest, err)
		} else if ext == "csv" {
			CSV(w, http.StatusOK, rows)
		} else if ext == "json" {
			JSON(w, http.StatusOK, rows)
		} else {
			ERROR(w, http.StatusNotFound, ErrEmpty)
		}
	}
}

// GET|POST /api/v1/query/<name>/explain
func (api *API) ExplainQuery(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	name, _ := parseName("name", r)
	if q := models.FindQuery(name); q == nil {
		ERROR(w, http.StatusNotFound, ErrEmpty)
	} else if q.Authorized(user) == false {
		ERROR(w, http.StatusUnauthorized, ErrEmpty)
	} else {
		params := parseParameters(r)
		if rows, err := q.Explain(params); err != nil {
			ERROR(w, http.StatusBadRequest, err)
		} else {
			JSON(w, http.StatusOK, rows)
		}
	}
}

// GET|POST /api/v1/query/<name>/<view_name>
func (api *API) RunView(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	name, _ := parseName("name", r)
	if query := models.FindQuery(name); query == nil {
		ERROR(w, http.StatusNotFound, ErrEmpty)
	} else if query.Authorized(user) == false {
		ERROR(w, http.StatusUnauthorized, ErrEmpty)
	} else {
		viewName, viewExt := parseName("view_name", r)
		view := query.View(viewName)
		if view == nil {
			ERROR(w, http.StatusNotFound, ErrEmpty)
			return
		}

		params := parseParameters(r)
		rows, err := query.Query(params)
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
			w.Header().Set("Content-Type", "image/svg+xml")
			w.WriteHeader(http.StatusOK)
			graph.Render(chart.SVG, w)
		} else {
			ERROR(w, http.StatusNotFound, ErrEmpty)
		}
	}
}
