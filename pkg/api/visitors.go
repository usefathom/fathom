package api

import (
	"net/http"

	"github.com/usefathom/fathom/pkg/count"
)

// URL: /api/visitors/count
var GetVisitorsCountHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	before, after := getRequestedPeriods(r)
	result := count.Visitors(before, after)
	return respond(w, envelope{Data: result})
})

// URL: /api/visitors/count/realtime
var GetVisitorsRealtimeCountHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	result := count.RealtimeVisitors()
	return respond(w, envelope{Data: result})
})

// URL: /api/visitors/count/group/:period
var GetVisitorsPeriodCountHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	before, after := getRequestedPeriods(r)
	results := count.VisitorsPerDay(before, after)
	return respond(w, envelope{Data: results})
})
