package count

import (
	"github.com/dannyvankooten/ana/datastore"
)

// TotalUniqueLanguages returns the total # of unique browser languages between two given timestamps
func TotalUniqueLanguages(before int64, after int64) int {
	var total int

	stmt, err := datastore.DB.Prepare(`
    SELECT
      IFNULL( SUM(t.count_unique), 0 )
    FROM total_browser_languages t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	checkError(err)
	defer stmt.Close()

	err = stmt.QueryRow(before, after).Scan(&total)
	checkError(err)

	return total
}

// Languages returns a point slice containing language data per language
func Languages(before int64, after int64, limit int) []Point {
	// TODO: Calculate total instead of requiring it as a parameter.
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
	total := TotalUniqueLanguages(before, after)
	points = calculatePointPercentages(points, total)

	return points
}

// CreateLanguageTotals aggregates screen data into daily totals
func CreateLanguageTotals(since string) {
	rows := queryTotalRows(`
    SELECT
      v.browser_language,
			COUNT(*) AS count,
      COUNT(DISTINCT(pv.visitor_id)) AS count_unique,
      DATE_FORMAT(pv.timestamp, "%Y-%m-%d") AS date_group
    FROM pageviews pv
    LEFT JOIN visitors v ON v.id = pv.visitor_id
    WHERE pv.timestamp > ?
    GROUP BY date_group, v.browser_language`, since)

	processTotalRows(rows, "total_browser_languages")
}
