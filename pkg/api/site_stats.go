package api

import (
	"net/http"

	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
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
	startDate, endDate := getRequestedDatePeriods(r)
	result, err := datastore.GetTotalSiteViews(startDate, endDate)
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
})

// URL: /api/stats/site/visitors
var GetSiteStatsVisitorsHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	startDate, endDate := getRequestedDatePeriods(r)
	result, err := datastore.GetTotalSiteVisitors(startDate, endDate)
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
})

// URL: /api/stats/site/duration
var GetSiteStatsDurationHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	startDate, endDate := getRequestedDatePeriods(r)
	result, err := datastore.GetAverageSiteDuration(startDate, endDate)
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
})

// URL: /api/stats/site/bounces
var GetSiteStatsBouncesHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	startDate, endDate := getRequestedDatePeriods(r)
	result, err := datastore.GetAverageSiteBounceRate(startDate, endDate)
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
})

// URL: /api/stats/site/realtime
var GetSiteStatsRealtimeHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	result, err := datastore.GetRealtimeVisitorCount()
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
})
