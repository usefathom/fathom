package count

import (
	"log"
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
)

// Browsers returns a point slice containing browser data per browser name
func Browsers(before int64, after int64, limit int) []Point {
	stmt, err := datastore.DB.Prepare(`
    SELECT
      t.value,
      SUM(t.count_unique) AS count
    FROM total_browser_names t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
    GROUP BY t.value
    ORDER BY count DESC
    LIMIT ?`)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query(before, after, limit)
	checkError(err)

	points := newPointSlice(rows)
	total, err := datastore.TotalUniqueBrowsers(before, after)
	checkError(err)

	points = calculatePointPercentages(points, total)

	return points
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
