package datastore

import "github.com/usefathom/fathom/pkg/models"

// TotalUniqueLanguages returns the total # of unique browser languages between two given timestamps
func TotalUniqueLanguages(before int64, after int64) (int, error) {
	var total int

	query := dbx.Rebind(`
    SELECT
      IFNULL( SUM(t.count_unique), 0 )
    FROM total_browser_languages t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)

	err := dbx.Get(&total, query, before, after)
	return total, err
}

func TotalsPerLanguage(before int64, after int64, limit int64) ([]*models.Point, error) {
	var results []*models.Point

	query := dbx.Rebind(`
		SELECT
	      t.value AS label,
	      SUM(t.count_unique) AS value
	   FROM total_browser_languages t
	   WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
	   GROUP BY label
	   ORDER BY value DESC
	   LIMIT ?`)

	err := dbx.Select(&results, query, before, after, limit)
	return results, err
}
