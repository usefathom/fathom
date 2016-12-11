package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/db"
  "github.com/dannyvankooten/ana/count"
  "encoding/json"
  "strings"
)

// URL: /api/referrers
var GetReferrersHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  before, after := getRequestedPeriods(r)

  // get total
  total := count.TotalVisitors(before, after)

  // get rows
  stmt, err := db.Conn.Prepare(`
    SELECT
    pv.referrer_url,
    COUNT(DISTINCT(pv.visitor_id)) AS count
    FROM pageviews pv
    WHERE UNIX_TIMESTAMP(pv.timestamp) <= ? AND UNIX_TIMESTAMP(pv.timestamp) >= ?
    AND pv.referrer_url IS NOT NULL
    AND pv.referrer_url != ""
    GROUP BY pv.referrer_url
    ORDER BY count DESC
    LIMIT ?`)
  checkError(err)
  defer stmt.Close()
  rows, err := stmt.Query(before, after, getRequestedLimit(r))
  checkError(err)
  defer rows.Close()

  results := make([]Datapoint, 0)
  for rows.Next() {
    var d Datapoint
    err = rows.Scan(&d.Label, &d.Count);
    d.Label = strings.Replace(d.Label, "http://", "", 1)
    d.Label = strings.Replace(d.Label, "https://", "", 1)
    d.Label = strings.TrimRight(d.Label, "/")
    checkError(err)

    d.Percentage = float64(d.Count) / total * 100
    results = append(results, d)
  }

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(results)
})
