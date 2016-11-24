package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/core"
  "encoding/json"
  "strconv"
)

// URL: /api/languages
var GetLanguagesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  stmt, err := core.DB.Prepare(`
    SELECT
    browser_language,
    COUNT(DISTINCT(ip_address)) AS count
    FROM visits
    WHERE timestamp >= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL ? DAY)
    GROUP BY browser_language
    ORDER BY count DESC`)
  checkError(err)
  defer stmt.Close()


  stmt2, err := core.DB.Prepare(`
    SELECT
    COUNT(DISTINCT(ip_address))
    FROM visits
    WHERE timestamp >= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL ? DAY)`)
  checkError(err)
  defer stmt2.Close()

  period, err := strconv.Atoi(r.URL.Query().Get("period"))
  if err != nil || period == 0 {
    period = defaultPeriod
  }

  var total float32
  stmt2.QueryRow(period).Scan(&total)

  rows, err := stmt.Query(period)
  checkError(err)

  type LanguageData struct {
    Count int
    Language string
    Percentage float32
  }

  results := make([]LanguageData, 0)
  defer rows.Close()
  for rows.Next() {
    var l LanguageData
    err = rows.Scan(&l.Language, &l.Count);
    checkError(err)

    l.Percentage = float32(l.Count) / total * 100
    results = append(results, l)
  }

  err = rows.Err();
  checkError(err)

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(results)
})
