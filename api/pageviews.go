package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/models"
  "github.com/dannyvankooten/ana/core"
  "encoding/json"
)

// URL: /api/pageviews
func GetPageviewsHandler(w http.ResponseWriter, r *http.Request) {
  stmt, err := core.DB.Prepare(`SELECT
      path,
      COUNT(DISTINCT(ip_address)) AS pageviews
    FROM visits
    GROUP BY path`)
  checkError(err)
  defer stmt.Close()

  rows, err := stmt.Query()
  checkError(err)

  results := make([]models.Pageview, 0)
  defer rows.Close()
  for rows.Next() {
    var p models.Pageview
    err = rows.Scan(&p.Path, &p.Count);
    checkError(err)
    results = append(results, p)
  }

  err = rows.Err();
  checkError(err)

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(results)
}
