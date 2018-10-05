package sqlstore

import (
	"database/sql"

	"github.com/usefathom/fathom/pkg/models"
)

// GetSites gets all sites in the database
func (db *sqlstore) GetSites() ([]*models.Site, error) {
	results := []*models.Site{}
	query := db.Rebind(`SELECT * FROM sites`)
	err := db.Select(&results, query)

	// don't err on no rows
	if err == sql.ErrNoRows {
		return results, nil
	}

	return results, err
}

func (db *sqlstore) GetSite(id int64) (*models.Site, error) {
	s := &models.Site{}
	query := db.Rebind("SELECT * FROM sites WHERE id = ?")
	err := db.Get(s, query, id)
	return s, mapError(err)
}

// SaveSite saves the website in the database (inserts or updates)
func (db *sqlstore) SaveSite(s *models.Site) error {
	if s.ID > 0 {
		return db.updateSite(s)
	}

	return db.insertSite(s)
}

// InsertSite saves a new site in the database
func (db *sqlstore) insertSite(s *models.Site) error {

	// Postgres does not support LastInsertID, so use a "... RETURNING" select query
	query := db.Rebind(`INSERT INTO sites(tracking_id, name) VALUES(?, ?)`)
	if db.Driver == POSTGRES {
		err := db.Get(&s.ID, query+" RETURNING id", s.TrackingID, s.Name)
		return err
	}

	// MySQL and SQLite do support LastInsertID, so use that
	r, err := db.Exec(query, s.TrackingID, s.Name)
	if err != nil {
		return err
	}

	s.ID, err = r.LastInsertId()
	return err

}

// UpdateSite updates an existing site in the database
func (db *sqlstore) updateSite(s *models.Site) error {
	query := db.Rebind(`UPDATE sites SET name = ? WHERE id = ?`)
	_, err := db.Exec(query, s.Name, s.ID)
	return err
}

// DeleteSite deletes the  given site in the database
func (db *sqlstore) DeleteSite(s *models.Site) error {
	query := db.Rebind(`DELETE FROM sites WHERE id = ?`)
	_, err := db.Exec(query, s.ID)
	return err
}
