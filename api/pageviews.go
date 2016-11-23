package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/models"
  "github.com/dannyvankooten/ana/core"
  "encoding/json"
)

// URL: /api/pageviews
var GetPageviewsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  stmt, err := core.DB.Prepare(`SELECT
      path,
      COUNT(ip_address) AS pageviews,
      COUNT(DISTINCT(ip_address)) AS pageviews_unique
    FROM visits
    WHERE timestamp >= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL ? DAY)
    GROUP BY path`)
  checkError(err)
  defer stmt.Close()

  period := r.URL.Query().Get("period")
  if period == "" {
    period = "1"
  }

  rows, err := stmt.Query(period)
  checkError(err)

  results := make([]models.Pageview, 0)
  defer rows.Close()
  for rows.Next() {
    var p models.Pageview
    err = rows.Scan(&p.Path, &p.Count, &p.CountUnique);
    checkError(err)
    results = append(results, p)
  }

  err = rows.Err();
  checkError(err)

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(results)
})


// URL: /api/pageviews/count/day
var GetPageviewsDayCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  stmt, err := core.DB.Prepare(`SELECT
    COUNT(*) AS count, DATE_FORMAT(timestamp, '%Y-%m-%d') AS date_group
    FROM visits
    WHERE timestamp >= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL ? DAY)
    GROUP BY date_group`)
  checkError(err)
  defer stmt.Close()

  period := r.URL.Query().Get("period")
  if period == "" {
    period = "1"
  }

  rows, err := stmt.Query(period)
  checkError(err)

  type Datapoint struct {
    Count int
    Label string
  }

  results := make([]Datapoint, 0)
  defer rows.Close()
  for rows.Next() {
    v := Datapoint{}
    err = rows.Scan(&v.Count, &v.Label);
    checkError(err)
    results = append(results, v)
  }

  err = rows.Err();
  checkError(err)

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(results)
})
