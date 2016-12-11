package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/count"
  "encoding/json"
)

// URL: /api/languages
var GetLanguagesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  before, after := getRequestedPeriods(r)

  // get total
  total := count.Visitors(before, after)

  results := count.Custom(`
    SELECT
    v.browser_language,
    COUNT(v.id) AS count
    FROM pageviews pv
    LEFT JOIN visitors v ON v.id = pv.visitor_id
    WHERE UNIX_TIMESTAMP(pv.timestamp) <= ? AND UNIX_TIMESTAMP(pv.timestamp) >= ?
    GROUP BY v.browser_language
    ORDER BY count DESC
    LIMIT ?`, before, after, getRequestedLimit(r), total,
  )

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(results)
})
