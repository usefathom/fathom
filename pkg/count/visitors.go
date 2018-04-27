package count

import (
	"log"
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
)

// CreateVisitorTotals aggregates visitor data into daily totals
func CreateVisitorTotals(since string) {
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	totals, err := datastore.VisitorCountPerDay(tomorrow, since)
	if err != nil {
		log.Fatal(err)
	}

	err = datastore.SaveVisitorTotals(totals)
	if err != nil {
		log.Fatal(err)
	}

}
