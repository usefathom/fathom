package datastore

import (
	"github.com/usefathom/fathom/pkg/models"
)

// TotalPageviews returns the total number of pageviews between the given timestamps
func TotalPageviews(before int64, after int64) (int, error) {
	var total int

	query := dbx.Rebind(`
		SELECT COALESCE(SUM(t.count), 0)
		FROM total_pageviews t 
		WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	err := dbx.Get(&total, query, before, after)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// TotalUniquePageviews returns the total number of unique pageviews between the given timestamps
func TotalUniquePageviews(before int64, after int64) (int, error) {
	var total int

	query := dbx.Rebind(`
		SELECT COALESCE(SUM(t.count_unique), 0)
		FROM total_pageviews t 
		WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	err := dbx.Get(&total, query, before, after)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// TotalPageviewsPerDay returns a slice of data points representing the number of pageviews per day
func TotalPageviewsPerDay(before int64, after int64) ([]*models.Total, error) {
	var results []*models.Total

	query := dbx.Rebind(`
		SELECT
	      CONCAT(p.scheme, "://", p.hostname, p.path) AS value,
		  COALESCE(SUM(t.count), 0) AS count,
		  COALESCE(SUM(t.count_unique), 0) AS count_unique,
	      DATE_FORMAT(t.date, '%Y-%m-%d') AS label
	    FROM total_pageviews t
	    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
	    GROUP BY label, p.hostname, p.path, p.scheme`)

	err := dbx.Select(&results, query, before, after)
	if err != nil {
		return results, err
	}

	return results, nil
}

// TotalPageviewsPerPage returns a set of pageview counts, grouped by page (hostname + path)
func TotalPageviewsPerPage(before int64, after int64, limit int64) ([]*models.Total, error) {
	var results []*models.Total
	query := dbx.Rebind(`
		SELECT
			CONCAT(p.scheme, "://", p.hostname, p.path) AS value,
			COALESCE(SUM(t.count), 0) AS count,
			COALESCE(SUM(t.count_unique), 0) AS count_unique
		FROM total_pageviews t
		LEFT JOIN pages p ON p.id = t.page_id
		WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
		GROUP BY p.hostname, p.path, p.scheme
		ORDER BY count DESC
		LIMIT ?`)
	err := dbx.Select(&results, query, before, after, limit)
	if err != nil {
		return results, err
	}

	return results, nil
}

// SavePageviewTotals saves the given totals in the connected database
// Differs slightly from the metric specific totals because of the normalized pages (to save storage)
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
		_, err = stmt.Exec(t.PageID, t.Count, t.CountUnique, t.Date, t.Count, t.CountUnique)
	}

	err = tx.Commit()
	return err
}
