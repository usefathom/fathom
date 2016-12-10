package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/db"
  "encoding/json"
)

// URL: /api/languages
var GetLanguagesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  before, after := getRequestedPeriods(r)

  stmt, err := db.Conn.Prepare(`
    SELECT
    COUNT(DISTINCT(pv.visitor_id))
    FROM pageviews pv
    LEFT JOIN visitors v ON v.id = pv.visitor_id
    WHERE UNIX_TIMESTAMP(pv.timestamp) <= ? AND UNIX_TIMESTAMP(pv.timestamp) >= ?
    `)
  checkError(err)
  defer stmt.Close()
  var total float32
  stmt.QueryRow(before, after).Scan(&total)

  stmt, err = db.Conn.Prepare(`
    SELECT
    v.browser_language,
    COUNT(v.id) AS count
    FROM pageviews pv
    LEFT JOIN visitors v ON v.id = pv.visitor_id
    WHERE UNIX_TIMESTAMP(pv.timestamp) <= ? AND UNIX_TIMESTAMP(pv.timestamp) >= ?
    GROUP BY v.browser_language
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
