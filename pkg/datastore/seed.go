package datastore

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/usefathom/fathom/pkg/models"
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
	SavePage(&homepage)

	contactPage := models.Page{
		Hostname: "wordpress.dev",
		Path:     "/contact/",
		Title:    "Contact",
	}
	SavePage(&contactPage)

	aboutPage := models.Page{
		Hostname: "wordpress.dev",
		Path:     "/about/",
		Title:    "About Me",
	}
	SavePage(&aboutPage)

	portfolioPage := models.Page{
		Hostname: "wordpress.dev",
		Path:     "/portfolio/",
		Title:    "Portfolio",
	}
	SavePage(&portfolioPage)

	pages = append(pages, homepage)
	pages = append(pages, homepage)
	pages = append(pages, contactPage)
	pages = append(pages, aboutPage)
	pages = append(pages, portfolioPage)
	return pages
}

// Seed inserts n random pageviews in the database.
func Seed(n int) {
	pages := seedPages()

	// insert X random hits
	for i := 0; i < n; i++ {

		// print a dot as progress indicator
		fmt.Print(".")
		date := randomDateBeforeNow()
		ipAddress := randomdata.IpV4Address()
		browserName := randSliceElement(browserNames)
		browserVersion := "54.0"
		deviceOS := "Linux"

		dummyUserAgent := browserName + browserVersion + deviceOS
		visitorKey := generateVisitorKey(date.Format("2006-01-02"), ipAddress, dummyUserAgent)
		visitor, err := GetVisitorByKey(visitorKey)

		if err != nil {
			// create or find visitor
			visitor := models.Visitor{
				IpAddress:        ipAddress,
				DeviceOS:         deviceOS,
				BrowserName:      browserName,
				BrowserVersion:   browserVersion,
				BrowserLanguage:  randSliceElement(browserLanguages),
				ScreenResolution: randSliceElement(screenResolutions),
				Country:          randomdata.Country(randomdata.TwoCharCountry),
			}
			err = SaveVisitor(&visitor)
		}

		// generate random timestamp
		timestamp := fmt.Sprintf("%s %d:%d:%d", date.Format("2006-01-02"), randInt(10, 24), randInt(10, 60), randInt(10, 60))

		pv := models.Pageview{
			VisitorID:       visitor.ID,
			ReferrerUrl:     "",
			ReferrerKeyword: "",
			Timestamp:       timestamp,
		}

		// insert between 1-6 pageviews for this visitor
		for j := 0; j <= randInt(1, 6); j++ {
			page := pages[randInt(0, len(pages))]
			pv.PageID = page.ID
			SavePageview(&pv)
		}
	}
}

func randomDate() time.Time {
	now := time.Now()
	month := months[randInt(0, len(months))]
	t := time.Date(randInt(now.Year()-1, now.Year()), month, randInt(1, 31), randInt(0, 23), randInt(0, 59), randInt(0, 59), 0, time.UTC)
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

// generateVisitorKey generates the "unique" visitor key from date, user agent + screen resolution
func generateVisitorKey(date string, ipAddress string, userAgent string) string {
	byteKey := md5.Sum([]byte(date + ipAddress + userAgent))
	return hex.EncodeToString(byteKey[:])
}
