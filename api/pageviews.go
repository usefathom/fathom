package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/models"
  "github.com/dannyvankooten/ana/core"
  "encoding/json"
)

// URL: /api/pageviews
var GetPageviewsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  period := getRequestedPeriod(r)

  stmt, err := core.DB.Prepare(`SELECT
      path,
      COUNT(ip_address) AS pageviews,
      COUNT(DISTINCT(ip_address)) AS pageviews_unique
    FROM visits
    WHERE timestamp >= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL ? DAY)
    GROUP BY path
    ORDER BY pageviews DESC`)
  checkError(err)
  defer stmt.Close()

  rows, err := stmt.Query(period)
  checkError(err)
  defer rows.Close()

  results := make([]models.Pageview, 0)
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
  period := getRequestedPeriod(r)

  stmt, err := core.DB.Prepare(`SELECT
    COUNT(*) AS count, DATE_FORMAT(timestamp, '%Y-%m-%d') AS date_group
    FROM visits
    WHERE timestamp >= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL ? DAY)
    GROUP BY date_group`)
  checkError(err)
  defer stmt.Close()

  rows, err := stmt.Query(period)
  checkError(err)
  defer rows.Close()

  results := make([]Datapoint, 0)
  for rows.Next() {
    v := Datapoint{}
    err = rows.Scan(&v.Count, &v.Label);
    checkError(err)
    results = append(results, v)
  }

  results = fillDatapoints(period, results)

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(results)
})
