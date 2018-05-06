package api

import (
	"github.com/usefathom/fathom/pkg/models"
	"net/http"
)

// URL: /api/stats/site
var GetSiteStatsHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	// before, after := getRequestedPeriods(r)
	// limit := getRequestedLimit(r)

	var results []*models.SiteStats
	return respond(w, envelope{Data: results})
})

// URL: /api/stats/site/pageviews
var GetSiteStatsPageviewsHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	// before, after := getRequestedPeriods(r)
	// limit := getRequestedLimit(r)
	var result int
	return respond(w, envelope{Data: result})
})

// URL: /api/stats/site/visitors
var GetSiteStatsVisitorsHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	// before, after := getRequestedPeriods(r)
	var result int
	return respond(w, envelope{Data: result})
})

// URL: /api/stats/site/avg-duration
var GetSiteStatsAvgDurationHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	// before, after := getRequestedPeriods(r)
	var result int
	return respond(w, envelope{Data: result})
})

// URL: /api/stats/site/avg-bounce
var GetSiteStatusAvgBounceHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	// before, after := getRequestedPeriods(r)
	var result int
	return respond(w, envelope{Data: result})
})
