package db

import (
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

var months = []time.Month {
  time.January,
  time.February,
  time.March,
  time.April,
  time.May,
  time.June,
  time.July,
  time.August,
  time.September,
  time.October,
  time.November,
  time.December,
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

func seedSite() models.Site {
  // get first site or create one
  var site models.Site
  Conn.QueryRow("SELECT id, url FROM sites LIMIT 1").Scan(&site.ID, &site.Url)

  if site.Url == "" {
    site.Url = "http://local.wordpress.dev/"
    site.Save(Conn)
  }

  return site
}

func seedPages(site models.Site) []models.Page {
  var pages = make([]models.Page, 0)

  homepage := models.Page{
    SiteID: site.ID,
    Path: "/",
    Title: "Homepage",
  }
  homepage.Save(Conn)

  contactPage := models.Page{
    SiteID: site.ID,
    Path: "/contact/",
    Title: "Contact",
  }
  contactPage.Save(Conn)

  aboutPage := models.Page{
    SiteID: site.ID,
    Path: "/about/",
    Title: "About Me",
  }
  aboutPage.Save(Conn)

  portfolioPage := models.Page{
    SiteID: site.ID,
    Path: "/portfolio/",
    Title: "Portfolio",
  }
  portfolioPage.Save(Conn)

  pages = append(pages, homepage)
  pages = append(pages, homepage)
  pages = append(pages, contactPage)
  pages = append(pages, aboutPage)
  pages = append(pages, portfolioPage)
  return pages
}

func Seed(n int) {

  site := seedSite()
  pages := seedPages(site)

  // prepare statement for inserting data
  stmt, err := Conn.Prepare(`INSERT INTO visits(
    page_id,
    browser_language,
    browser_name,
    browser_version,
    country,
    ip_address,
    referrer_url,
    screen_resolution,
    timestamp
    ) VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ? )`)
  if err != nil {
      log.Fatal(err)
  }
  defer stmt.Close()

  // insert X random hits
  log.Printf("Inserting %d visits", n)
  for i := 0; i < n; i++ {

    // print a dot as progress indicator
    fmt.Print(".")

    // generate random timestamp
    date := randomDateBeforeNow();
    timestamp := fmt.Sprintf("%s %d:%d:%d", date.Format("2006-01-02"), randInt(10, 24), randInt(10, 60), randInt(10, 60))

    visit := models.Visit{
      IpAddress: randomdata.IpV4Address(),
      BrowserName: randSliceElement(browserNames),
      BrowserVersion: "54.0.2840.100",
      BrowserLanguage: randSliceElement(browserLanguages),
      ScreenResolution: randSliceElement(screenResolutions),
      Country: randomdata.Country(randomdata.TwoCharCountry),
      ReferrerUrl: "",
      Timestamp: timestamp,
    }

    // insert between 1-4 pageviews for this visitor
    for j := 0; j < randInt(1, 4); j++ {
      page := pages[randInt(0, len(pages))]
      visit.PageID = page.ID

      _, err = stmt.Exec(
        visit.PageID,
        visit.BrowserLanguage,
        visit.BrowserName,
        visit.BrowserVersion,
        visit.Country,
        visit.IpAddress,
        visit.ReferrerUrl,
        visit.ScreenResolution,
        visit.Timestamp,
      )
      if err != nil {
        log.Fatal(err)
      }
    }

  }
}

func randomDate() time.Time {
  now := time.Now()
  month := months[randInt(0, len(months))]
  t := time.Date(now.Year(), month, randInt(1,31), randInt(0,23), randInt(0,59), randInt(0,59), 0, time.UTC)
  return t
}

func randomDateBeforeNow() time.Time {
  now := time.Now()
  date := randomDate()
  for( date.After(now) ) {
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
