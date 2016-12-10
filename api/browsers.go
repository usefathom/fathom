package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/db"
  "encoding/json"
)

// URL: /api/browsers
var GetBrowsersHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  before, after := getRequestedPeriods(r)

  // get total
  stmt, err := db.Conn.Prepare(`
    SELECT
    COUNT(DISTINCT(pv.visitor_id))
    FROM pageviews pv
    WHERE UNIX_TIMESTAMP(pv.timestamp) <= ? AND UNIX_TIMESTAMP(pv.timestamp) >= ?
  `)
  checkError(err)
  defer stmt.Close()
  var total float32
  stmt.QueryRow(before, after).Scan(&total)

  // get rows
  stmt, err = db.Conn.Prepare(`
    SELECT
    v.browser_name,
    COUNT(DISTINCT(pv.visitor_id)) AS count
    FROM pageviews pv
    LEFT JOIN visitors v ON v.id = pv.visitor_id
    WHERE UNIX_TIMESTAMP(pv.timestamp) <= ? AND UNIX_TIMESTAMP(pv.timestamp) >= ? AND v.browser_name IS NOT NULL
    GROUP BY v.browser_name
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
