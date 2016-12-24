package count

import (
	"github.com/dannyvankooten/ana/db"
)

// Screens returns a point slice containing screen data per size
func Screens(before int64, after int64, limit int, total float64) []Point {
	// TODO: Calculate total instead of requiring it as a parameter.
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

	return newPointSlice(rows, total)
}

// CreateScreenArchives aggregates screen data into daily totals
func CreateScreenArchives() {
	rows := queryTotalRows(`
    SELECT
      v.screen_resolution,
			COUNT(*) AS count,
      COUNT(DISTINCT(pv.visitor_id)) AS count_unique,
      DATE_FORMAT(pv.timestamp, "%Y-%m-%d") AS date_group
    FROM pageviews pv
    LEFT JOIN visitors v ON v.id = pv.visitor_id
    WHERE NOT EXISTS(
      SELECT t.id
      FROM total_screens t
      WHERE t.date = DATE_FORMAT(pv.timestamp, "%Y-%m-%d")
    )
    GROUP BY date_group, v.screen_resolution`)

	processTotalRows(rows, "total_screens")
}
