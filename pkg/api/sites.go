package api

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/usefathom/fathom/pkg/models"
)

// seed rand pkg on program init
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

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
	var s *models.Site
	vars := mux.Vars(r)
	sid, ok := vars["id"]
	if ok {
		id, err := strconv.ParseInt(sid, 10, 64)
		if err != nil {
			return err
		}

		s, err = api.database.GetSite(id)
		if err != nil {
			return err
		}
	} else {
		s = &models.Site{
			TrackingID: generateTrackingID(),
		}
	}

	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		return err
	}

	log.Printf("Site tracking ID: %s\n", s.TrackingID)
	if err := api.database.SaveSite(s); err != nil {
		return err
	}

	// TODO: If we just created the first site, add existing data (with site_id = 0) to the site we just created

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

func generateTrackingID() string {
	return randomString(5)
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //a=65 and z = 65+25
	}

	return string(bytes)
}
