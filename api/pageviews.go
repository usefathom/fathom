package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/models"
  "github.com/dannyvankooten/ana/db"
  "encoding/json"
  "github.com/gorilla/mux"
  "time"
)

// URL: /api/pageviews
var GetPageviewsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  before, after := getRequestedPeriods(r)
  stmt, err := db.Conn.Prepare(`SELECT
      s.url,
      p.path,
      COUNT(v.ip_address) AS pageviews,
      COUNT(DISTINCT(v.ip_address)) AS pageviews_unique
    FROM visits v
    LEFT JOIN pages p ON v.page_id = p.id
    LEFT JOIN sites s ON p.site_id = s.id
    WHERE UNIX_TIMESTAMP(v.timestamp) <= ? AND UNIX_TIMESTAMP(timestamp) >= ?
    GROUP BY p.path, s.url
    ORDER BY pageviews DESC
    LIMIT ?`)
  checkError(err)
  defer stmt.Close()

  rows, err := stmt.Query(before, after, defaultLimit)
  checkError(err)
  defer rows.Close()

  results := make([]models.Pageview, 0)
  for rows.Next() {
    var p models.Pageview
    err = rows.Scan(&p.Url, &p.Path, &p.Count, &p.CountUnique);
    checkError(err)
    results = append(results, p)
  }

  err = rows.Err();
  checkError(err)

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(results)
})


// URL: /api/pageviews/count
var GetPageviewsCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  before, after := getRequestedPeriods(r)
  stmt, err := db.Conn.Prepare(`SELECT COUNT(*) FROM visits WHERE UNIX_TIMESTAMP(timestamp) <= ? AND UNIX_TIMESTAMP(timestamp) >= ?`)
  checkError(err)
  defer stmt.Close()

  var result int
  stmt.QueryRow(before, after).Scan(&result)

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(result)
})

// URL: /api/pageviews/group/day
var GetPageviewsPeriodCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  period := vars["period"]
  formats := map[string]string {
    "day": "%Y-%m-%d",
    "month": "%Y-%m",
  }
  before, after := getRequestedPeriods(r)
  stmt, err := db.Conn.Prepare(`SELECT
    COUNT(*) AS count, DATE_FORMAT(timestamp, ?) AS date_group
    FROM visits
    WHERE UNIX_TIMESTAMP(timestamp) <= ? AND UNIX_TIMESTAMP(timestamp) >= ?
    GROUP BY date_group`)
  checkError(err)
  defer stmt.Close()

  rows, err := stmt.Query(formats[period], before, after)
  checkError(err)
  defer rows.Close()

  results := make([]Datapoint, 0)
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
