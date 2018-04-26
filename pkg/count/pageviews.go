package count

import (
	"github.com/usefathom/fathom/pkg/datastore"
	"log"
	"time"
)

// CreatePageviewTotals aggregates pageview data for each page into daily totals
func CreatePageviewTotals(since string) {
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-02-01")
	totals, err := datastore.PageviewCountPerPageAndDay(tomorrow, since)
	if err != nil {
		log.Fatal(err)
	}

	err = datastore.SavePageviewTotals(totals)
	if err != nil {
		log.Fatal(err)
	}
}
