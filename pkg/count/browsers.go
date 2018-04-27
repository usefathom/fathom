package count

import (
	"log"
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
)

// Browsers returns a point slice containing browser data per browser name
func Browsers(before int64, after int64, limit int64) ([]*models.Point, error) {
	points, err := datastore.TotalsPerBrowser(before, after, limit)
	if err != nil {
		return nil, err
	}

	total, err := datastore.TotalUniqueBrowsers(before, after)
	if err != nil {
		return nil, err
	}

	points = calculatePointPercentages(points, total)
	return points, nil
}

// CreateBrowserTotals aggregates screen data into daily totals
func CreateBrowserTotals(since string) {
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	totals, err := datastore.BrowserCountPerDay(tomorrow, since)
	if err != nil {
		log.Fatal(err)
	}

	err = datastore.SaveTotals("browser_names", totals)
	if err != nil {
		log.Fatal(err)
	}
}
