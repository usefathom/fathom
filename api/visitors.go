package api

import (
	"encoding/json"
	"net/http"

	"github.com/dannyvankooten/ana/count"
)

// URL: /api/visitors/count
var GetVisitorsCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)
	result := count.Visitors(before, after)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
})

// URL: /api/visitors/count/realtime
var GetVisitorsRealtimeCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	result := count.RealtimeVisitors()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
})

// URL: /api/visitors/count/group/:period
var GetVisitorsPeriodCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)
	results := count.VisitorsPerDay(before, after)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
})
