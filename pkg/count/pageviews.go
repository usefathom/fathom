package count

import "github.com/usefathom/fathom/pkg/datastore"

// CreatePageviewTotals aggregates pageview data for each page into daily totals
func CreatePageviewTotals(since string) {
	stmt, err := datastore.DB.Prepare(`SELECT
      pv.page_id,
      COUNT(*) AS count,
			COUNT(DISTINCT(pv.visitor_id)) AS count_unique,
			DATE_FORMAT(pv.timestamp, "%Y-%m-%d") AS date_group
    FROM pageviews pv
    WHERE pv.timestamp > ?
    GROUP BY pv.page_id, date_group`)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query(since)
	checkError(err)
	defer rows.Close()

	datastore.DB.Begin()

	datastore.DB.Exec("START TRANSACTION")
	for rows.Next() {
		var t Total
		err = rows.Scan(&t.PageID, &t.Count, &t.CountUnique, &t.Date)
		checkError(err)

		stmt, err := datastore.DB.Prepare(`INSERT INTO total_pageviews(
	    page_id,
	    count,
			count_unique,
	    date
	    ) VALUES( ?, ?, ?, ? ) ON DUPLICATE KEY UPDATE count = ?, count_unique = ?`)
		checkError(err)
		defer stmt.Close()

		_, err = stmt.Exec(
			t.PageID,
			t.Count,
			t.CountUnique,
			t.Date,
			t.Count,
			t.CountUnique,
		)
		checkError(err)
	}
	datastore.DB.Exec("COMMIT")
}
