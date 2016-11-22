package api

import (
  "net/http"
  "github.com/dannyvankooten/ana/models"
  "github.com/dannyvankooten/ana/core"
  "encoding/json"
  "strconv"
)

// URL: /api/visits
func GetVisitsHandler(w http.ResponseWriter, r *http.Request) {
  stmt, err := core.DB.Prepare(`SELECT
    id,
    COALESCE(browser_name, '') AS browser_name,
    COALESCE(browser_language, '') AS browser_language,
    COALESCE(device_os, '') AS device_os,
    ip_address,
    path,
    COALESCE(screen_resolution, '') AS screen_resolution,
    timestamp
    FROM visits
    ORDER BY timestamp DESC
    LIMIT ?`)

  checkError(err)
  defer stmt.Close()

  limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
  if limit == 0 {
    limit = 10
  }

  rows, err := stmt.Query(&limit)
  checkError(err)

  results := make([]models.Visit, 0)
  defer rows.Close()
  for rows.Next() {
    var v models.Visit
    err = rows.Scan(&v.ID, &v.BrowserName, &v.BrowserLanguage, &v.DeviceOS, &v.IpAddress, &v.Path, &v.ScreenResolution, &v.Timestamp);
    checkError(err)
    results = append(results, v)
  }

  err = rows.Err();
  checkError(err)

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(results)
}

// URL: /api/visits/count/realtime
func GetVisitsRealtimeCount(w http.ResponseWriter, r *http.Request) {
  row := core.DB.QueryRow(`SELECT COUNT(DISTINCT(ip_address)) FROM visits WHERE timestamp >= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL 3 HOUR_MINUTE)`)
  var result int
  row.Scan(&result)

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(result)
}

// URL: /api/visits/count/day
func GetVisitsDayCountHandler(w http.ResponseWriter, r *http.Request) {
  stmt, err := core.DB.Prepare(`SELECT
    COUNT(*) AS count, DATE_FORMAT(timestamp, '%Y-%m-%d') AS date_group
    FROM visits
    GROUP BY date_group`)
  checkError(err)
  defer stmt.Close()

  rows, err := stmt.Query()
  checkError(err)

  type Datapoint struct {
    Count int
    Label string
  }

  results := make([]Datapoint, 0)
  defer rows.Close()
  for rows.Next() {
    v := Datapoint{}
    err = rows.Scan(&v.Count, &v.Label);
    checkError(err)
    results = append(results, v)
  }

  err = rows.Err();
  checkError(err)

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(results)
}
