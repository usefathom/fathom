package count

import (
	"github.com/dannyvankooten/ana/db"
)

// Pageviews returns the total number of pageviews between the given timestamps
func Pageviews(before int64, after int64) float64 {
	// get total
	stmt, err := db.Conn.Prepare(`
    SELECT
    SUM(t.count)
    FROM total_pageviews t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	checkError(err)
	defer stmt.Close()
	var total float64
	stmt.QueryRow(before, after).Scan(&total)
	return total
}

// PageviewsPerDay returns a slice of data points representing the number of pageviews per day
func PageviewsPerDay(before int64, after int64) []Point {
	stmt, err := db.Conn.Prepare(`SELECT
      SUM(t.count) AS count,
      DATE_FORMAT(t.date, '%Y-%m-%d') AS date_group
    FROM total_pageviews t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
    GROUP BY date_group`)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query(before, after)
	checkError(err)
	defer rows.Close()

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

// CreatePageviewArchives aggregates pageview data for each page into daily totals
func CreatePageviewArchives() {
	stmt, err := db.Conn.Prepare(`SELECT
      pv.page_id,
      COUNT(*) AS count,
			COUNT(DISTINCT(pv.visitor_id)) AS count_unique,
			DATE_FORMAT(pv.timestamp, "%Y-%m-%d") AS date_group
    FROM pageviews pv
    WHERE NOT EXISTS (
			SELECT t.id
			FROM total_pageviews t
			WHERE t.date = DATE_FORMAT(pv.timestamp, "%Y-%m-%d") AND t.page_id = pv.page_id
		)
    GROUP BY pv.page_id, date_group`)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkError(err)
	defer rows.Close()

	db.Conn.Exec("START TRANSACTION")
	for rows.Next() {
		var t Total
		err = rows.Scan(&t.PageID, &t.Count, &t.CountUnique, &t.Date)
		checkError(err)

		stmt, err := db.Conn.Prepare(`INSERT INTO total_pageviews(
	    page_id,
	    count,
			count_unique,
	    date
	    ) VALUES( ?, ?, ?, ? )`)
		checkError(err)
		defer stmt.Close()

		_, err = stmt.Exec(
			t.PageID,
			t.Count,
			t.CountUnique,
			t.Date,
		)
		checkError(err)
	}
	db.Conn.Exec("COMMIT")
}
