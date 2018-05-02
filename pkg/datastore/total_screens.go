package datastore

import "github.com/usefathom/fathom/pkg/models"

// TotalScreens returns the total # of screens between two given timestamps
func TotalScreens(before int64, after int64) (int, error) {
	var total int

	query := dbx.Rebind(`
       SELECT
         COALESCE(SUM(t.count), 0)
       FROM total_screens t
       WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	err := dbx.Get(&total, query, before, after)
	return total, err
}

// TotalUniqueScreens returns the total # of unique screens between two given timestamps
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

func TotalsPerScreen(before int64, after int64, limit int64) ([]*models.Total, error) {
	var results []*models.Total

	query := dbx.Rebind(`
      SELECT
        t.value AS value,
        COALESCE(SUM(t.count), 0) AS count,
        COALESCE(SUM(t.count_unique), 0) AS count_unique
      FROM total_screens t
      WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
      GROUP BY t.value
      ORDER BY count DESC
      LIMIT ?`)

	err := dbx.Select(&results, query, before, after, limit)
	return results, err
}
