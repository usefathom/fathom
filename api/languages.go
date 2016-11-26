package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/core"
  "encoding/json"
)

// URL: /api/languages
var GetLanguagesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  before, after := getRequestedPeriods(r)

  stmt, err := core.DB.Prepare(`
    SELECT
    COUNT(DISTINCT(ip_address))
    FROM visits
    WHERE UNIX_TIMESTAMP(timestamp) <= ? AND UNIX_TIMESTAMP(timestamp) >= ?
    `)
  checkError(err)
  defer stmt.Close()
  var total float32
  stmt.QueryRow(before, after).Scan(&total)

  stmt, err = core.DB.Prepare(`
    SELECT
    browser_language,
    COUNT(DISTINCT(ip_address)) AS count
    FROM visits
    WHERE UNIX_TIMESTAMP(timestamp) <= ? AND UNIX_TIMESTAMP(timestamp) >= ?
    GROUP BY browser_language
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
