package api

import (
	"net/http"

	"github.com/usefathom/fathom/pkg/datastore"
)

// URL: /api/stats/page
var GetPageStatsHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := datastore.GetAggregatedPageStats(params.StartDate, params.EndDate, params.Limit)
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
})

// URL: /api/stats/page/pageviews
var GetPageStatsPageviewsHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := datastore.GetAggregatedPageStatsPageviews(params.StartDate, params.EndDate)
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
})
