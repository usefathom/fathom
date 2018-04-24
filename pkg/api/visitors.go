package api

import (
	"net/http"

	"github.com/usefathom/fathom/pkg/count"
)

// URL: /api/visitors/count
var GetVisitorsCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)
	result := count.Visitors(before, after)
	respond(w, envelope{Data: result})
})

// URL: /api/visitors/count/realtime
var GetVisitorsRealtimeCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	result := count.RealtimeVisitors()
	respond(w, envelope{Data: result})
})

// URL: /api/visitors/count/group/:period
var GetVisitorsPeriodCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)
	results := count.VisitorsPerDay(before, after)
	respond(w, envelope{Data: results})
})
