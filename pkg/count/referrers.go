package count

import (
	"log"
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
)

// Referrers returns a point slice containing browser data per browser name
func Referrers(before int64, after int64, limit int) []Point {
	stmt, err := datastore.DB.Prepare(`
    SELECT
      t.value,
      SUM(t.count) AS count
    FROM total_referrers t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
    GROUP BY t.value
    ORDER BY count DESC
    LIMIT ?`)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query(before, after, limit)
	checkError(err)

	points := newPointSlice(rows)
	total, err := datastore.TotalReferrers(before, after)
	checkError(err)

	points = calculatePointPercentages(points, total)

	return points
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
