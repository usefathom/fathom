package api

import (
	"github.com/usefathom/fathom/pkg/models"
	"net/http"
)

// URL: /api/stats/page
var GetPageStatsHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	//before, after := getRequestedPeriods(r)
	//limit := getRequestedLimit(r)
	var result []*models.PageStats
	return respond(w, envelope{Data: result})
})
