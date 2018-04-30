package datastore

import "github.com/usefathom/fathom/pkg/models"

// TotalVisitors returns the number of total visitors between the given timestamps
func TotalVisitors(before int64, after int64) (int, error) {
	var total int

	query := dbx.Rebind(`
    SELECT COALESCE(SUM(t.count), 0)
    FROM total_visitors t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	err := dbx.Get(&total, query, before, after)
	return total, err
}

// TotalVisitorsPerDay returns a point slice containing visitor data per day
func TotalVisitorsPerDay(before int64, after int64) ([]*models.Point, error) {
	var results []*models.Point

	query := dbx.Rebind(`SELECT
      COALESCE(SUM(t.count), 0) AS value,
      DATE_FORMAT(t.date, '%Y-%m-%d') AS label
    FROM total_visitors t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
    GROUP BY label`)

	err := dbx.Select(&results, query, before, after)
	return results, err
}

func SaveVisitorTotals(totals []*models.Total) error {
	tx, err := dbx.Begin()
	if err != nil {
		return nil
	}

	query := dbx.Rebind(`INSERT INTO total_visitors( count, date ) VALUES( ?, ? ) ON DUPLICATE KEY UPDATE count = ?`)
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	for _, t := range totals {
		_, err = stmt.Exec(t.Count, t.Date, t.Count)
	}

	err = tx.Commit()
	return err
}
