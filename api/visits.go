package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/models"
  "github.com/dannyvankooten/ana/db"
  "encoding/json"
  "github.com/gorilla/mux"
  "time"
)

// URL: /api/visits
var GetVisitsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  stmt, err := db.Conn.Prepare(`SELECT
    id,
    COALESCE(browser_name, '') AS browser_name,
    COALESCE(browser_language, '') AS browser_language,
    ip_address,
    COALESCE(screen_resolution, '') AS screen_resolution,
    timestamp
    FROM visits
    ORDER BY timestamp DESC
    LIMIT ?`)

  checkError(err)
  defer stmt.Close()

  limit := getRequestedLimit(r)
  rows, err := stmt.Query(&limit)
  checkError(err)

  results := make([]models.Visit, 0)
  defer rows.Close()
  for rows.Next() {
    var v models.Visit
    err = rows.Scan(&v.ID, &v.BrowserName, &v.BrowserLanguage, &v.IpAddress, &v.ScreenResolution, &v.Timestamp);
    checkError(err)
    results = append(results, v)
  }

  err = rows.Err();
  checkError(err)

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(results)
})

// URL: /api/visits/count
var GetVisitsCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  before, after := getRequestedPeriods(r)
  stmt, err := db.Conn.Prepare(`SELECT COUNT(DISTINCT(ip_address)) FROM visits WHERE UNIX_TIMESTAMP(timestamp) <= ? AND UNIX_TIMESTAMP(timestamp) >= ?`)

  checkError(err)
  defer stmt.Close()

  var result int
  stmt.QueryRow(before, after).Scan(&result)

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(result)
})

// URL: /api/visits/count/realtime
var GetVisitsRealtimeCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  var result int
  db.Conn.QueryRow(`SELECT COUNT(DISTINCT(ip_address)) FROM visits WHERE timestamp >= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL 3 HOUR_MINUTE) AND timestamp <= CURRENT_TIMESTAMP`).Scan(&result)

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(result)
})

// URL: /api/visits/count/group/:period
var GetVisitsPeriodCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  period := vars["period"]
  formats := map[string]string {
    "day": "%Y-%m-%d",
    "month": "%Y-%m",
  }

  stmt, err := db.Conn.Prepare(`SELECT
    COUNT(DISTINCT(ip_address)) AS count, DATE_FORMAT(timestamp, ?) AS date_group
    FROM visits
    WHERE UNIX_TIMESTAMP(timestamp) <= ? AND UNIX_TIMESTAMP(timestamp) >= ?
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
