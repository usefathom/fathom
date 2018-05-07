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
	params := GetRequestParams(r)
	result, err := datastore.GetTotalSiteViews(params.StartDate, params.EndDate)
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
})

// URL: /api/stats/site/visitors
var GetSiteStatsVisitorsHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := datastore.GetTotalSiteVisitors(params.StartDate, params.EndDate)
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
})

// URL: /api/stats/site/duration
var GetSiteStatsDurationHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := datastore.GetAverageSiteDuration(params.StartDate, params.EndDate)
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
})

// URL: /api/stats/site/bounces
var GetSiteStatsBouncesHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := datastore.GetAverageSiteBounceRate(params.StartDate, params.EndDate)
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
