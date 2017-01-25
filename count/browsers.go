package count

import (
	"github.com/dannyvankooten/ana/datastore"
)

// TotalUniqueBrowsers returns the total # of unique browsers between two given timestamps
func TotalUniqueBrowsers(before int64, after int64) int {
	var total int

	stmt, err := datastore.DB.Prepare(`
    SELECT
      IFNULL( SUM(t.count_unique), 0 )
    FROM total_browser_names t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	checkError(err)
	defer stmt.Close()

	err = stmt.QueryRow(before, after).Scan(&total)
	checkError(err)

	return total
}

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
	total := TotalUniqueBrowsers(before, after)
	points = calculatePointPercentages(points, total)

	return points
}

// CreateBrowserTotals aggregates screen data into daily totals
func CreateBrowserTotals(since string) {
	rows := queryTotalRows(`
    SELECT
      v.browser_name,
			COUNT(*) AS count,
      COUNT(DISTINCT(pv.visitor_id)) AS count_unique,
      DATE_FORMAT(pv.timestamp, "%Y-%m-%d") AS date_group
    FROM pageviews pv
    LEFT JOIN visitors v ON v.id = pv.visitor_id
    WHERE pv.timestamp > ?
    GROUP BY date_group, v.browser_name`, since)

	processTotalRows(rows, "total_browser_names")
}
