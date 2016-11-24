package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/core"
  "encoding/json"
)

// URL: /api/screen-resolutions
var GetScreenResolutionsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  period := getRequestedPeriod(r)

  // get total
  stmt, err := core.DB.Prepare(`
    SELECT
    COUNT(DISTINCT(ip_address))
    FROM visits
    WHERE timestamp >= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL ? DAY)`)
  checkError(err)
  defer stmt.Close()
  var total float32
  stmt.QueryRow(period).Scan(&total)

  // get rows
  stmt, err = core.DB.Prepare(`
    SELECT
    screen_resolution,
    COUNT(DISTINCT(ip_address)) AS count
    FROM visits
    WHERE timestamp >= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL ? DAY)
    GROUP BY screen_resolution
    ORDER BY count DESC`)
  checkError(err)
  defer stmt.Close()
  rows, err := stmt.Query(period)
  checkError(err)
  defer rows.Close()

  results := make([]Datapoint, 0)
  for rows.Next() {
    var d Datapoint
    err = rows.Scan(&d.Label, &d.Count);
    checkError(err)

    d.Percentage = float32(d.Count) / total * 100
    results = append(results, d)
  }

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(results)
})
