package api

import (
	"net/http"

	"github.com/usefathom/fathom/pkg/count"
	"github.com/usefathom/fathom/pkg/datastore"
)

// URL: /api/pageviews
var GetPageviewsHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	before, after := getRequestedPeriods(r)
	limit := getRequestedLimit(r)
	results, err := count.Pageviews(before, after, limit)
	if err != nil {
		return err
	}

	return respond(w, envelope{Data: results})
})

// URL: /api/pageviews/count
var GetPageviewsCountHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	before, after := getRequestedPeriods(r)
	result, err := datastore.TotalPageviews(before, after)
	if err != nil {
		return err
	}

	return respond(w, envelope{Data: result})
})

// URL: /api/pageviews/group/day
var GetPageviewsPeriodCountHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	before, after := getRequestedPeriods(r)
	results, err := datastore.TotalPageviewsPerDay(before, after)
	if err != nil {
		return err
	}

	return respond(w, envelope{Data: results})
})
