package db

import (
  "github.com/dannyvankooten/ana/core"
  "github.com/dannyvankooten/ana/models"
  "log"
  "time"
  "math/rand"
  "fmt"
  "github.com/Pallinder/go-randomdata"
)

var browserNames = []string {
  "Chrome",
  "Chrome",
  "Firefox",
  "Safari",
  "Safari",
  "Internet Explorer",
}

var paths = []string {
  "/",
  "/", // we need this to weigh more.
  "/contact",
  "/about",
  "/checkout",
}

var browserLanguages = []string {
  "en-US",
  "en-US",
  "nl-NL",
  "fr-FR",
  "de-DE",
  "es-ES",
}

var screenResolutions = []string {
  "2560x1440",
  "1920x1080",
  "1920x1080",
  "360x640",
}

func Seed(n int) {

  // prepare statement for inserting data
  stmt, err := core.DB.Prepare(`INSERT INTO visits(
    browser_language,
    browser_name,
    browser_version,
    country,
    device_os,
    ip_address,
    path,
    referrer_url,
    screen_resolution,
    timestamp
    ) VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )`)
  if err != nil {
      log.Fatal(err)
  }
  defer stmt.Close()

  // insert X random hits
  for i := 0; i < n; i++ {

    // generate random timestamp
    date := randomDateBeforeNow();
    timestamp := fmt.Sprintf("%s %d:%d:%d", date.Format("2006-01-02"), randInt(10, 24), randInt(10, 60), randInt(10, 60))

    visit := models.Visit{
      Path: randSliceElement(paths),
      IpAddress: randomdata.IpV4Address(),
      DeviceOS: "Linux x86_64",
      BrowserName: randSliceElement(browserNames),
      BrowserVersion: "54.0.2840.100",
      BrowserLanguage: randSliceElement(browserLanguages),
      ScreenResolution: randSliceElement(screenResolutions),
      Country: randomdata.Country(randomdata.TwoCharCountry),
      ReferrerUrl: "",
      Timestamp: timestamp,
    }

    _, err = stmt.Exec(
      visit.BrowserLanguage,
      visit.BrowserName,
      visit.BrowserVersion,
      visit.Country,
      visit.DeviceOS,
      visit.IpAddress,
      visit.Path,
      visit.ReferrerUrl,
      visit.ScreenResolution,
      visit.Timestamp,
    )
    if err != nil {
      log.Fatal(err)
    }

  }
}

func randomDate() time.Time {
  date, _ := time.Parse("Monday 2 Jan 2006", randomdata.FullDate())
  return date
}

func randomDateBeforeNow() time.Time {
  date := randomDate()
  for( date.After(time.Now()) ) {
    date = randomDate()
  }

  return date
}

func randSliceElement(slice []string) string {
  return slice[randInt(0, len(slice))]
}

func randInt(min int, max int) int {
    return min + rand.Intn(max-min)
}
