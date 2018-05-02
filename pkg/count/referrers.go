package count

import (
	"log"
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
)

// Referrers returns a point slice containing browser data per browser name
func Referrers(before int64, after int64, limit int64) ([]*models.Total, error) {
	points, err := datastore.TotalsPerReferrer(before, after, limit)
	if err != nil {
		return nil, err
	}

	total, err := datastore.TotalReferrers(before, after)
	if err != nil {
		return nil, err
	}

	points = calculatePercentagesOfTotal(points, total)
	return points, nil
}

// CreateReferrerTotals aggregates screen data into daily totals
func CreateReferrerTotals(since string) {
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	totals, err := datastore.ReferrerCountPerDay(tomorrow, since)
	if err != nil {
		log.Fatal(err)
	}

	err = datastore.SaveTotals("referrers", totals)
	if err != nil {
		log.Fatal(err)
	}
}
