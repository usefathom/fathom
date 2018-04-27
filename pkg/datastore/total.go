package datastore

import (
	"fmt"

	"github.com/usefathom/fathom/pkg/models"
)

func SaveTotals(metric string, totals []*models.Total) error {
	query := dbx.Rebind(fmt.Sprintf(`
		INSERT INTO total_%s( value, count, count_unique, date) 
		VALUES( ?, ?, ?, ? ) ON DUPLICATE KEY UPDATE count = ?, count_unique = ?
	`, metric))

	tx, err := dbx.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, t := range totals {
		result, err := stmt.Exec(t.Value, t.Count, t.CountUnique, t.Date, t.Count, t.CountUnique)
		if err != nil {
			return err
		}

		t.ID, _ = result.LastInsertId()
	}

	err = tx.Commit()
	return err
}
