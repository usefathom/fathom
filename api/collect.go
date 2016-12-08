package api

import (
  "net/http"
  "log"
  "github.com/mssola/user_agent"
  "github.com/dannyvankooten/ana/models"
  "github.com/dannyvankooten/ana/db"
  "encoding/base64"
)

func CollectHandler(w http.ResponseWriter, r *http.Request) {
  ua := user_agent.New(r.UserAgent())

  // abort if this is a bot.
  if ua.Bot() {
    return
  }

  // prepare statement for inserting data
  stmt, err := db.Conn.Prepare(`INSERT INTO visits(
    ip_address,
    page_id,
    referrer_url,
    browser_language,
    browser_name,
    browser_version,
    screen_resolution
    ) VALUES( ?, ?, ?, ?, ?, ?, ?, ? )`)
  if err != nil {
      log.Fatal(err.Error())
  }
  defer stmt.Close()

  // TODO: Mask IP Address
  // TODO: Query DB to determine whether visitor is returning
  ipAddress := r.RemoteAddr
  headerForwardedFor := r.Header.Get("X-Forwarded-For")
  if( headerForwardedFor != "" ) {
    ipAddress = headerForwardedFor
  }

  // TODO: Query Path

  q := r.URL.Query()
  visit := models.Visit{
    IpAddress: ipAddress,
    ReferrerUrl: q.Get("r"),
    BrowserLanguage: q.Get("l"),
    ScreenResolution: q.Get("sr"),
  }

  // add browser details
  visit.BrowserName, visit.BrowserVersion = ua.Browser()

  _, err = stmt.Exec(
    visit.IpAddress,
    visit.ReferrerUrl,
    visit.BrowserLanguage,
    visit.BrowserName,
    visit.BrowserVersion,
    visit.ScreenResolution,
  )
  if err != nil {
    log.Fatal(err)
  }

  // don't you cache this
  w.Header().Set("Content-Type", "image/gif")
  w.Header().Set("Expires", "Mon, 01 Jan 1990 00:00:00 GMT")
  w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
  w.Header().Set("Pragma", "no-cache")
  w.WriteHeader(http.StatusOK)

  // 1x1 px transparent GIF
  b, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7")
  w.Write(b)
}
