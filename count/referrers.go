package count

import (
	"github.com/dannyvankooten/ana/db"
)

// Referrers returns a point slice containing browser data per browser name
func Referrers(before int64, after int64, limit int, total float64) []Point {
	// TODO: Calculate total instead of requiring it as a parameter.
	stmt, err := db.Conn.Prepare(`
    SELECT
      t.value,
      SUM(t.count_unique) AS count
    FROM total_referrers t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
    GROUP BY t.value
    ORDER BY count DESC
    LIMIT ?`)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query(before, after, limit)
	checkError(err)

	total := Visitors(before, after)

	return newPointSlice(rows, total)
}

// CreateReferrerArchives aggregates screen data into daily totals
func CreateReferrerArchives() {
	rows := queryTotalRows(`
    SELECT
      pv.referrer_url,
			COUNT(*) AS count,
      COUNT(DISTINCT(pv.visitor_id)) AS count_unique,
      DATE_FORMAT(pv.timestamp, "%Y-%m-%d") AS date_group
    FROM pageviews pv
    WHERE pv.referrer_url IS NOT NULL
    AND pv.referrer_url != ''
    AND NOT EXISTS(
      SELECT t.id
      FROM total_referrers t
      WHERE t.date = DATE_FORMAT(pv.timestamp, "%Y-%m-%d")
    )
    GROUP BY date_group, pv.referrer_url`)

	processTotalRows(rows, "total_referrers")
}
