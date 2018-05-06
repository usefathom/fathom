package api

import (
	"github.com/usefathom/fathom/pkg/models"
	"net/http"
)

// URL: /api/stats/referrer
var GetReferrerStatsHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	// before, after := getRequestedPeriods(r)
	// limit := getRequestedLimit(r)
	var result []*models.ReferrerStats
	return respond(w, envelope{Data: result})
})
