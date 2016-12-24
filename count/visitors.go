package count

import (
	"github.com/dannyvankooten/ana/db"
)

// TODO: Convert to total_visitors table.

// Visitors returns the number of total visitors between the given timestamps
func Visitors(before int64, after int64) float64 {
	// get total
	stmt, err := db.Conn.Prepare(`
    SELECT
    SUM(t.count)
    FROM total_visitors t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	checkError(err)
	defer stmt.Close()
	var total float64
	stmt.QueryRow(before, after).Scan(&total)
	return total
}

// VisitorsPerDay returns a point slice containing visitor data per day
func VisitorsPerDay(before int64, after int64) []Point {
	stmt, err := db.Conn.Prepare(`SELECT
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

	results = fill(after, before, results)

	return results
}

// CreateVisitorArchives aggregates visitor data into daily totals
func CreateVisitorArchives() {
	stmt, err := db.Conn.Prepare(`
    SELECT
      COUNT(DISTINCT(pv.visitor_id)) AS count,
      DATE_FORMAT(pv.timestamp, "%Y-%m-%d") AS date_group
    FROM pageviews pv
    WHERE NOT EXISTS(
      SELECT t.id
      FROM total_visitors t
      WHERE t.date = DATE_FORMAT(pv.timestamp, "%Y-%m-%d")
    )
    GROUP BY date_group`)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkError(err)
	defer rows.Close()

	db.Conn.Exec("START TRANSACTION")
	for rows.Next() {
		var t Total
		err = rows.Scan(&t.Count, &t.Date)
		checkError(err)

		stmt, err := db.Conn.Prepare(`INSERT INTO total_visitors(
	    count,
	    date
	    ) VALUES( ?, ? )`)
		checkError(err)
		defer stmt.Close()

		_, err = stmt.Exec(
			t.Count,
			t.Date,
		)
		checkError(err)
	}
	db.Conn.Exec("COMMIT")
}
