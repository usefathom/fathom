package count

import (
	"log"
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
)

// CreatePageviewTotals aggregates pageview data for each page into daily totals
func CreatePageviewTotals(since string) {
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	totals, err := datastore.PageviewCountPerPageAndDay(tomorrow, since)
	if err != nil {
		log.Fatal(err)
	}

	err = datastore.SavePageviewTotals(totals)
	if err != nil {
		log.Fatal(err)
	}
}
