package datastore

import "github.com/usefathom/fathom/pkg/models"

// TotalUniqueScreens returns the total # of screens between two given timestamps
func TotalUniqueScreens(before int64, after int64) (int, error) {
	var total int

	query := dbx.Rebind(`
       SELECT
         COALESCE(SUM(t.count_unique), 0)
       FROM total_screens t
       WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	err := dbx.Get(&total, query, before, after)
	return total, err
}

func TotalsPerScreen(before int64, after int64, limit int64) ([]*models.Point, error) {
	var results []*models.Point

	query := dbx.Rebind(`
      SELECT
        t.value AS label,
        COALESCE(SUM(t.count_unique), 0) AS value
      FROM total_screens t
      WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
      GROUP BY t.value
      ORDER BY value DESC
      LIMIT ?`)

	err := dbx.Select(&results, query, before, after, limit)
	return results, err
}
