package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/usefathom/fathom/pkg/models"
)

// GET /api/sites
func (api *API) GetSitesHandler(w http.ResponseWriter, r *http.Request) error {
	result, err := api.database.GetSites()
	if err != nil {
		return err
	}
	return respond(w, http.StatusOK, envelope{Data: result})
}

// POST /api/sites
// POST /api/sites/{id}
func (api *API) SaveSiteHandler(w http.ResponseWriter, r *http.Request) error {
	s := &models.Site{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		return err
	}

	if err := api.database.SaveSite(s); err != nil {
		return err
	}

	return respond(w, http.StatusOK, envelope{Data: s})
}

// DELETE /api/sites/{id}
func (api *API) DeleteSiteHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return err
	}

	if err := api.database.DeleteSite(&models.Site{ID: id}); err != nil {
		return err
	}

	return respond(w, http.StatusOK, envelope{Data: true})
}
