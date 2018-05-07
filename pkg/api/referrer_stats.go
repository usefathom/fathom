package api

import (
	"net/http"

	"github.com/usefathom/fathom/pkg/datastore"
)

// URL: /api/stats/referrer
var GetReferrerStatsHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := datastore.GetAggregatedReferrerStats(params.StartDate, params.EndDate, params.Limit)
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
})
