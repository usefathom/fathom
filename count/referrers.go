package count

import (
	"github.com/dannyvankooten/ana/datastore"
)

// TotalReferrers returns the total # of referrers between two given timestamps
func TotalReferrers(before int64, after int64) int {
	var total int

	stmt, err := datastore.DB.Prepare(`
    SELECT
      IFNULL( SUM(t.count), 0 )
    FROM total_referrers t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	checkError(err)
	defer stmt.Close()

	err = stmt.QueryRow(before, after).Scan(&total)
	checkError(err)

	return total
}

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
	total := TotalReferrers(before, after)
	points = calculatePointPercentages(points, total)

	return points
}

// CreateReferrerTotals aggregates screen data into daily totals
func CreateReferrerTotals(since string) {
	rows := queryTotalRows(`
    SELECT
      pv.referrer_url,
			COUNT(*) AS count,
      COUNT(DISTINCT(pv.visitor_id)) AS count_unique,
      DATE_FORMAT(pv.timestamp, "%Y-%m-%d") AS date_group
    FROM pageviews pv
    WHERE pv.referrer_url IS NOT NULL
    AND pv.referrer_url != ''
    AND pv.timestamp > ?
    GROUP BY date_group, pv.referrer_url`, since)

	processTotalRows(rows, "total_referrers")
}
