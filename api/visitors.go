package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/db"
  "github.com/dannyvankooten/ana/count"
  "encoding/json"
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
  var result int
  db.Conn.QueryRow(`
    SELECT COUNT(DISTINCT(pv.visitor_id))
    FROM pageviews pv
    WHERE pv.timestamp >= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL 3 HOUR_MINUTE) AND pv.timestamp <= CURRENT_TIMESTAMP`).Scan(&result)
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
