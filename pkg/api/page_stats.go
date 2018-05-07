package api

import (
	"net/http"

	"github.com/usefathom/fathom/pkg/datastore"
)

// URL: /api/stats/page
var GetPageStatsHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	startDate, endDate := getRequestedDatePeriods(r)
	limit := getRequestedLimit(r)
	result, err := datastore.GetAggregatedPageStats(startDate, endDate, limit)
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
})
