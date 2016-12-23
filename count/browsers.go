package count

import (
	"github.com/dannyvankooten/ana/db"
)

// Browsers returns a point slice containing browser data per browser name
func Browsers(before int64, after int64, limit int, total float64) []Point {
	// TODO: Calculate total instead of requiring it as a parameter.
	stmt, err := db.Conn.Prepare(`
    SELECT
      a.value,
      SUM(a.count) AS count
    FROM archive a
    WHERE a.metric = 'browsers' AND UNIX_TIMESTAMP(a.date) <= ? AND UNIX_TIMESTAMP(a.date) >= ?
    GROUP BY a.value
    ORDER BY count DESC
    LIMIT ?`)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query(before, after, limit)
	checkError(err)

	return newPointSlice(rows, total)
}

// CreateBrowserArchives aggregates screen data into daily totals
func CreateBrowserArchives() {
	stmt, err := db.Conn.Prepare(`
    SELECT
      v.browser_name, 
      COUNT(DISTINCT(pv.visitor_id)) AS count,
      DATE_FORMAT(pv.timestamp, "%Y-%m-%d") AS date_group
    FROM pageviews pv
    LEFT JOIN visitors v ON v.id = pv.visitor_id
    WHERE NOT EXISTS(
      SELECT a.id
      FROM archive a
      WHERE a.metric = 'browsers' AND a.date = DATE_FORMAT(pv.timestamp, "%Y-%m-%d")
    )
    GROUP BY date_group, v.browser_name`)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkError(err)
	defer rows.Close()

	db.Conn.Exec("START TRANSACTION")
	for rows.Next() {
		a := Archive{
			Metric: "browsers",
		}
		err = rows.Scan(&a.Value, &a.Count, &a.Date)
		checkError(err)
		a.Save(db.Conn)
	}
	db.Conn.Exec("COMMIT")
}
