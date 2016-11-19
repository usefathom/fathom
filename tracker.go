package ana

import (
  "net/http"
  "log"
  "time"
  "github.com/mssola/user_agent"
)

func collectHandler(w http.ResponseWriter, r *http.Request) {
  log.Printf("[%s] %s %s (%s)\n", time.Now(), r.Method, r.RequestURI, r.UserAgent())

  ua := user_agent.New(r.UserAgent())
  if ua.Bot() {
    return
  }

  // prepare statement for inserting data
  stmt, err := db.Prepare(`INSERT INTO visits(
    path,
    ip_address,
    referrer_url,
    browser_language,
    browser_name,
    browser_version,
    device_os,
    screen_resolution
    ) VALUES( ?, ?, ?, ?, ?, ?, ?, ? )`)
  if err != nil {
      log.Fatal(err.Error())
  }
  defer stmt.Close()

  // TODO: Mask IP Address
  // TODO: Query DB to determine whether visitor is returning

  q := r.URL.Query()
  visit := Visit{
    Path: q.Get("p"),
    IpAddress: r.RemoteAddr,
    ReferrerUrl: q.Get("r"),
    BrowserLanguage: q.Get("l"),
    ScreenResolution: q.Get("sr"),
  }

  // add browser details
  visit.BrowserName, visit.BrowserVersion = ua.Browser()

  // add device details
  visit.DeviceOS = ua.OS()

  log.Printf("%+v\n", visit)

  _, err = stmt.Exec(
    visit.IpAddress,
    visit.Path,
    visit.ReferrerUrl,
    visit.BrowserLanguage,
    visit.BrowserName,
    visit.BrowserVersion,
    visit.DeviceOS,
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
  w.Header().Set("Status", "200")
}

func main() {
    db = connectToDatabase()

    http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
      http.ServeFile(w, r, r.URL.Path[1:])
    })
    http.HandleFunc("/collect", collectHandler)
    http.ListenAndServe(":8080", nil)
}
