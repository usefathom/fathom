package datastore

import (
	"github.com/usefathom/fathom/pkg/models"
)

// TotalPageviews returns the total number of pageviews between the given timestamps
func TotalPageviews(before int64, after int64) (int, error) {
	var total int

	query := dbx.Rebind(`SELECT IFNULL( SUM(t.count), 0 ) FROM total_pageviews t WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	err := dbx.Get(&total, query, before, after)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// TotalPageviewsPerDay returns a slice of data points representing the number of pageviews per day
func TotalPageviewsPerDay(before int64, after int64) ([]*models.Point, error) {
	var results []*models.Point

	query := dbx.Rebind(`
		SELECT
	      SUM(t.count) AS value,
	      DATE_FORMAT(t.date, '%Y-%m-%d') AS label
	    FROM total_pageviews t
	    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
	    GROUP BY label`)

	err := dbx.Select(&results, query, before, after)
	if err != nil {
		return results, err
	}

	return results, nil
}

// TotalPageviewsPerPage returns a set of pageview counts, grouped by page (hostname + path)
func TotalPageviewsPerPage(before int64, after int64, limit int64) ([]*models.PageviewCount, error) {
	var results []*models.PageviewCount
	query := dbx.Rebind(`
		SELECT
			p.hostname,
			p.path,
			SUM(t.count) AS count,
			SUM(t.count_unique) AS countunique
		FROM total_pageviews t
		LEFT JOIN pages p ON p.id = t.page_id
		WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
		GROUP BY p.path, p.hostname
		ORDER BY count DESC
		LIMIT ?`)
	err := dbx.Select(&results, query, before, after, limit)
	if err != nil {
		return results, err
	}

	return results, nil
}

func SavePageviewTotals(totals []*models.Total) error {
	tx, err := dbx.Begin()
	if err != nil {
		return nil
	}

	query := dbx.Rebind(`INSERT INTO total_pageviews( page_id, count, count_unique, date ) VALUES( ?, ?, ?, ? ) ON DUPLICATE KEY UPDATE count = ?, count_unique = ?`)
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	for _, t := range totals {
		_, err = stmt.Exec(
			t.PageID,
			t.Count,
			t.CountUnique,
			t.Date,
			t.Count,
			t.CountUnique,
		)
	}

	err = tx.Commit()
	return err
}
