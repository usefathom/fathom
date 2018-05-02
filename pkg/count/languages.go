package count

import (
	"log"
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
)

// Languages returns a point slice containing language data per language
func Languages(before int64, after int64, limit int64) ([]*models.Total, error) {
	points, err := datastore.TotalsPerLanguage(before, after, limit)
	if err != nil {
		return nil, err
	}

	total, err := datastore.TotalLanguages(before, after)
	if err != nil {
		return nil, err
	}

	points = calculatePercentagesOfTotal(points, total)
	return points, nil
}

// CreateLanguageTotals aggregates screen data into daily totals
func CreateLanguageTotals(since string) {
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	totals, err := datastore.LanguageCountPerDay(tomorrow, since)
	if err != nil {
		log.Fatal(err)
	}

	err = datastore.SaveTotals("browser_languages", totals)
	if err != nil {
		log.Fatal(err)
	}
}
