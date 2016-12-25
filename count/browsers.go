package count

import (
	"github.com/dannyvankooten/ana/db"
)

// Browsers returns a point slice containing browser data per browser name
func Browsers(before int64, after int64, limit int) []Point {
	stmt, err := db.Conn.Prepare(`
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

	return newPointSlice(rows)
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
