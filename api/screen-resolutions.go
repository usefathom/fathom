package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/core"
  "encoding/json"
)

// URL: /api/screen-resolutions
var GetScreenResolutionsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  before, after := getRequestedPeriods(r)

  // get total
  stmt, err := core.DB.Prepare(`
    SELECT
    COUNT(DISTINCT(ip_address))
    FROM visits
    WHERE UNIX_TIMESTAMP(timestamp) <= ? AND UNIX_TIMESTAMP(timestamp) >= ?`)
  checkError(err)
  defer stmt.Close()
  var total float32
  stmt.QueryRow(before, after).Scan(&total)

  // get rows
  stmt, err = core.DB.Prepare(`
    SELECT
    screen_resolution,
    COUNT(DISTINCT(ip_address)) AS count
    FROM visits
    WHERE UNIX_TIMESTAMP(timestamp) <= ? AND UNIX_TIMESTAMP(timestamp) >= ?
    GROUP BY screen_resolution
    ORDER BY count DESC
    LIMIT ?`)
  checkError(err)
  defer stmt.Close()
  rows, err := stmt.Query(before, after, defaultLimit)
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
