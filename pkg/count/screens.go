package count

import (
	"log"
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
)

// Screens returns a point slice containing screen data per size
func Screens(before int64, after int64, limit int64) ([]*models.Total, error) {
	points, err := datastore.TotalsPerScreen(before, after, limit)
	if err != nil {
		return nil, err
	}

	total, err := datastore.TotalScreens(before, after)
	if err != nil {
		return nil, err
	}

	points = calculatePercentagesOfTotal(points, total)
	return points, nil
}

// CreateScreenTotals aggregates screen data into daily totals
func CreateScreenTotals(since string) {
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	totals, err := datastore.ScreenCountPerDay(tomorrow, since)
	if err != nil {
		log.Fatal(err)
	}

	err = datastore.SaveTotals("screens", totals)
	if err != nil {
		log.Fatal(err)
	}
}
