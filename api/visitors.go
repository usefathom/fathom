package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/db"
  "encoding/json"
  "github.com/gorilla/mux"
  "time"
)

// URL: /api/visitors/count
var GetVisitorsCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  before, after := getRequestedPeriods(r)
  stmt, err := db.Conn.Prepare(`SELECT COUNT(DISTINCT(visitor_id)) FROM pageviews WHERE UNIX_TIMESTAMP(timestamp) <= ? AND UNIX_TIMESTAMP(timestamp) >= ?`)

  checkError(err)
  defer stmt.Close()

  var result int
  stmt.QueryRow(before, after).Scan(&result)
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
  vars := mux.Vars(r)
  period := vars["period"]
  formats := map[string]string {
    "day": "%Y-%m-%d",
    "month": "%Y-%m",
  }

  stmt, err := db.Conn.Prepare(`SELECT
    COUNT(DISTINCT(visitor_id)) AS count, DATE_FORMAT(timestamp, ?) AS date_group
    FROM pageviews pv
    WHERE UNIX_TIMESTAMP(pv.timestamp) <= ? AND UNIX_TIMESTAMP(pv.timestamp) >= ?
    GROUP BY date_group`)
  checkError(err)
  defer stmt.Close()

  before, after := getRequestedPeriods(r)
  rows, err := stmt.Query(formats[period], before, after)
  checkError(err)

  results := make([]Datapoint, 0)
  defer rows.Close()
  for rows.Next() {
    v := Datapoint{}
    err = rows.Scan(&v.Count, &v.Label);
    checkError(err)
    results = append(results, v)
  }

  d := time.Hour * 24;
  if period == "month" {
    d = d * 30
  }
  results = fillDatapoints(before, after, d, results)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(results)
})
