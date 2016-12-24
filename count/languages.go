package count

import (
	"github.com/dannyvankooten/ana/db"
)

// Languages returns a point slice containing language data per language
func Languages(before int64, after int64, limit int) []Point {
	// TODO: Calculate total instead of requiring it as a parameter.
	stmt, err := db.Conn.Prepare(`
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

	total := Visitors(before, after)

	return newPointSlice(rows, total)
}

// CreateLanguageArchives aggregates screen data into daily totals
func CreateLanguageArchives() {
	rows := queryTotalRows(`
    SELECT
      v.browser_language,
			COUNT(*) AS count,
      COUNT(DISTINCT(pv.visitor_id)) AS count_unique,
      DATE_FORMAT(pv.timestamp, "%Y-%m-%d") AS date_group
    FROM pageviews pv
    LEFT JOIN visitors v ON v.id = pv.visitor_id
    WHERE NOT EXISTS(
      SELECT t.id
      FROM total_browser_languages t
      WHERE t.date = DATE_FORMAT(pv.timestamp, "%Y-%m-%d")
    )
    GROUP BY date_group, v.browser_language`)

	processTotalRows(rows, "total_browser_languages")
}
