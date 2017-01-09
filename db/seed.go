package db

import (
	"fmt"
	"math/rand"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/dannyvankooten/ana/models"
)

var browserNames = []string{
	"Chrome",
	"Chrome",
	"Firefox",
	"Safari",
	"Safari",
	"Internet Explorer",
}

var months = []time.Month{
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

var browserLanguages = []string{
	"en-US",
	"en-US",
	"nl-NL",
	"fr-FR",
	"de-DE",
	"es-ES",
}

var screenResolutions = []string{
	"2560x1440",
	"1920x1080",
	"1920x1080",
	"360x640",
}

func seedPages() []models.Page {
	var pages = make([]models.Page, 0)

	homepage := models.Page{
		Hostname: "wordpress.dev",
		Path:     "/",
		Title:    "Homepage",
	}
	homepage.Save(Conn)

	contactPage := models.Page{
		Hostname: "wordpress.dev",
		Path:     "/contact/",
		Title:    "Contact",
	}
	contactPage.Save(Conn)

	aboutPage := models.Page{
		Hostname: "wordpress.dev",
		Path:     "/about/",
		Title:    "About Me",
	}
	aboutPage.Save(Conn)

	portfolioPage := models.Page{
		Hostname: "wordpress.dev",
		Path:     "/portfolio/",
		Title:    "Portfolio",
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
	pages := seedPages()

	stmtVisitor, _ := Conn.Prepare("SELECT v.id FROM visitors v WHERE v.visitor_key = ? LIMIT 1")
	defer stmtVisitor.Close()

	// insert X random hits
	for i := 0; i < n; i++ {

		// print a dot as progress indicator
		fmt.Print(".")

		// create or find visitor
		visitor := models.Visitor{
			IpAddress:        randomdata.IpV4Address(),
			DeviceOS:         "Linux",
			BrowserName:      randSliceElement(browserNames),
			BrowserVersion:   "54.0",
			BrowserLanguage:  randSliceElement(browserLanguages),
			ScreenResolution: randSliceElement(screenResolutions),
			Country:          randomdata.Country(randomdata.TwoCharCountry),
		}
		visitor.Key = visitor.GenerateKey()

		err := stmtVisitor.QueryRow(visitor.Key).Scan(&visitor.ID)
		if err != nil {
			visitor.Save(Conn)
		}

		// generate random timestamp
		date := randomDateBeforeNow()
		timestamp := fmt.Sprintf("%s %d:%d:%d", date.Format("2006-01-02"), randInt(10, 24), randInt(10, 60), randInt(10, 60))

		pv := models.Pageview{
			VisitorID:       visitor.ID,
			ReferrerUrl:     "",
			ReferrerKeyword: "",
			Timestamp:       timestamp,
		}

		Conn.Exec("START TRANSACTION")

		// insert between 1-4 pageviews for this visitor
		for j := 0; j <= randInt(1, 4); j++ {
			page := pages[randInt(0, len(pages))]
			pv.PageID = page.ID
			pv.Save(Conn)
		}

		Conn.Exec("COMMIT")
	}
}

func randomDate() time.Time {
	now := time.Now()
	month := months[randInt(0, len(months))]
	t := time.Date(now.Year(), month, randInt(1, 31), randInt(0, 23), randInt(0, 59), randInt(0, 59), 0, time.UTC)
	return t
}

func randomDateBeforeNow() time.Time {
	now := time.Now()
	date := randomDate()
	for date.After(now) {
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
