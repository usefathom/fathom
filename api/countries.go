package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/count"
  "encoding/json"
)

// URL: /api/countries
var GetCountriesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  before, after := getRequestedPeriods(r)

  // get total
  total := count.Visitors(before, after)

  // get rows
  results := count.Custom(`
    SELECT
    v.country,
    COUNT(DISTINCT(pv.visitor_id)) AS count
    FROM pageviews pv
    LEFT JOIN visitors v ON v.id = pv.visitor_id
    WHERE UNIX_TIMESTAMP(pv.timestamp) <= ? AND UNIX_TIMESTAMP(pv.timestamp) >= ? AND v.country IS NOT NULL
    GROUP BY v.country
    ORDER BY count DESC
    LIMIT ?`, before, after, getRequestedLimit(r), total)

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(results)
})
