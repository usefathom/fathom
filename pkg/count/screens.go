package count

import (
	"log"
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
)

// Screens returns a point slice containing screen data per size
func Screens(before int64, after int64, limit int) []Point {
	stmt, err := datastore.DB.Prepare(`
    SELECT
      t.value,
      SUM(t.count_unique) AS count
    FROM total_screens t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
    GROUP BY t.value
    ORDER BY count DESC
    LIMIT ?`)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query(before, after, limit)
	checkError(err)

	points := newPointSlice(rows)
	total, err := datastore.TotalUniqueScreens(before, after)
	checkError(err)

	points = calculatePointPercentages(points, total)

	return points
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
