package count

import (
	"log"
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
)

// CreateBouncesTotals aggregates pageview data for each page into daily totals
func CreateBouncesTotals(since string) {
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	totals, err := datastore.BouncesCountPerPageAndDay(tomorrow, since)
	if err != nil {
		log.Fatal(err)
	}

	err = datastore.SavePageTotals("bounced", totals)
	if err != nil {
		log.Fatal(err)
	}
}
