package api

import (
  "net/http"
  "log"
  "github.com/dannyvankooten/ana/models"
  "github.com/dannyvankooten/ana/core"
  "encoding/json"
)

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
    FROM visits`)
  if err != nil {
    log.Fatal(err.Error())
  }
  defer stmt.Close()


  rows, err := stmt.Query()
  if err != nil {
    log.Fatal(err.Error())
  }

  results := make([]models.Visit, 0)
  defer rows.Close()
  for rows.Next() {
    var v models.Visit
    if err := rows.Scan(&v.ID, &v.BrowserName, &v.BrowserLanguage, &v.DeviceOS, &v.IpAddress, &v.Path, &v.ScreenResolution, &v.Timestamp); err != nil {
      log.Fatal(err)
    }

    results = append(results, v)
    }
    if err := rows.Err(); err != nil {
      log.Fatal(err)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
  }
