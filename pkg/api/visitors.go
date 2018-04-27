package api

import (
	"net/http"

	"github.com/usefathom/fathom/pkg/datastore"
)

// URL: /api/visitors/count
var GetVisitorsCountHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	before, after := getRequestedPeriods(r)
	result, err := datastore.TotalVisitors(before, after)
	if err != nil {
		return err
	}

	return respond(w, envelope{Data: result})
})

// URL: /api/visitors/count/realtime
var GetVisitorsRealtimeCountHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	result, err := datastore.RealtimeVisitors()
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
})

// URL: /api/visitors/count/group/:period
var GetVisitorsPeriodCountHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	before, after := getRequestedPeriods(r)
	results, err := datastore.TotalVisitorsPerDay(before, after)
	if err != nil {
		return err
	}

	return respond(w, envelope{Data: results})
})
