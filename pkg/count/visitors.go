package count

import (
	"github.com/usefathom/fathom/pkg/datastore"
)

// RealtimeVisitors returns the total number of visitors in the last 3 minutes
func RealtimeVisitors() int {
	var result int
	datastore.DB.QueryRow(`
		SELECT COUNT(DISTINCT(pv.visitor_id))
		FROM pageviews pv
		WHERE pv.timestamp >= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL 3 HOUR_MINUTE) AND pv.timestamp <= CURRENT_TIMESTAMP`).Scan(&result)
	return result
}

// Visitors returns the number of total visitors between the given timestamps
func Visitors(before int64, after int64) int {
	// get total
	stmt, err := datastore.DB.Prepare(`
    SELECT
    SUM(t.count)
    FROM total_visitors t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	checkError(err)
	defer stmt.Close()
	var total int
	stmt.QueryRow(before, after).Scan(&total)
	return total
}

// VisitorsPerDay returns a point slice containing visitor data per day
func VisitorsPerDay(before int64, after int64) []Point {
	stmt, err := datastore.DB.Prepare(`SELECT
      SUM(t.count) AS count,
      DATE_FORMAT(t.date, '%Y-%m-%d') AS date_group
    FROM total_visitors t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
    GROUP BY date_group`)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query(before, after)
	checkError(err)

	var results []Point
	defer rows.Close()
	for rows.Next() {
		p := Point{}
		err = rows.Scan(&p.Value, &p.Label)
		checkError(err)
		results = append(results, p)
	}

	return results
}

// CreateVisitorTotals aggregates visitor data into daily totals
func CreateVisitorTotals(since string) {
	stmt, err := datastore.DB.Prepare(`
    SELECT
      COUNT(DISTINCT(pv.visitor_id)) AS count,
      DATE_FORMAT(pv.timestamp, "%Y-%m-%d") AS date_group
    FROM pageviews pv
    WHERE pv.timestamp > ?
    GROUP BY date_group`)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query(since)
	checkError(err)
	defer rows.Close()

	datastore.DB.Exec("START TRANSACTION")
	for rows.Next() {
		var t Total
		err = rows.Scan(&t.Count, &t.Date)
		checkError(err)

		stmt, err := datastore.DB.Prepare(`INSERT INTO total_visitors(
	    count,
	    date
	    ) VALUES( ?, ? ) ON DUPLICATE KEY UPDATE count = ?`)
		checkError(err)
		defer stmt.Close()

		_, err = stmt.Exec(
			t.Count,
			t.Date,
			t.Count,
		)
		checkError(err)
	}
	datastore.DB.Exec("COMMIT")
}
