package count

import (
	"log"
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
)

// Languages returns a point slice containing language data per language
func Languages(before int64, after int64, limit int) []Point {
	stmt, err := datastore.DB.Prepare(`
    SELECT
      t.value,
      SUM(t.count_unique) AS count
    FROM total_browser_languages t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
    GROUP BY t.value
    ORDER BY count DESC
    LIMIT ?`)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query(before, after, limit)
	checkError(err)

	points := newPointSlice(rows)
	total, err := datastore.TotalUniqueLanguages(before, after)
	checkError(err)

	points = calculatePointPercentages(points, total)

	return points
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
