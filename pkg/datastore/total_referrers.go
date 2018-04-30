package datastore

import "github.com/usefathom/fathom/pkg/models"

// TotalReferrers returns the total # of referrers between two given timestamps
func TotalReferrers(before int64, after int64) (int, error) {
	var total int

	query := dbx.Rebind(`
       SELECT
         COALESCE(SUM(t.count), 0)
       FROM total_referrers t
       WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)

	err := dbx.Get(&total, query, before, after)
	return total, err
}

func TotalsPerReferrer(before int64, after int64, limit int64) ([]*models.Point, error) {
	var results []*models.Point

	query := dbx.Rebind(`
      SELECT
         t.value AS label,
         COALESCE(SUM(t.count), 0) AS value
      FROM total_referrers t
      WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
      GROUP BY label
      ORDER BY value DESC
      LIMIT ?`)

	err := dbx.Select(&results, query, before, after, limit)
	return results, err
}
