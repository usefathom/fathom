package count

import (
	"github.com/dannyvankooten/ana/db"
)

// TotalUniqueScreens returns the total # of screens between two given timestamps
func TotalUniqueScreens(before int64, after int64) int {
	var total int

	stmt, err := db.Conn.Prepare(`
    SELECT
      SUM(t.count_unique)
    FROM total_screens t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	checkError(err)
	defer stmt.Close()

	err = stmt.QueryRow(before, after).Scan(&total)
	checkError(err)

	return total
}

// Screens returns a point slice containing screen data per size
func Screens(before int64, after int64, limit int) []Point {
	stmt, err := db.Conn.Prepare(`
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
	total := TotalUniqueScreens(before, after)
	points = calculatePointPercentages(points, total)

	return points
}

// CreateScreenTotals aggregates screen data into daily totals
func CreateScreenTotals(since string) {
	rows := queryTotalRows(`
    SELECT
      v.screen_resolution,
			COUNT(*) AS count,
      COUNT(DISTINCT(pv.visitor_id)) AS count_unique,
      DATE_FORMAT(pv.timestamp, "%Y-%m-%d") AS date_group
    FROM pageviews pv
    LEFT JOIN visitors v ON v.id = pv.visitor_id
    WHERE pv.timestamp > ?
    GROUP BY date_group, v.screen_resolution`, since)

	processTotalRows(rows, "total_screens")
}
