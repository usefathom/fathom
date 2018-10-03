package sqlstore

import (
	"github.com/usefathom/fathom/pkg/models"
)

// GetSites gets all sites in the database
func (db *sqlstore) GetSites() ([]*models.Site, error) {
	var results []*models.Site
	query := db.Rebind(`SELECT * FROM sites`)
	err := db.Select(&results, query)
	return results, err
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
	query := db.Rebind(`INSERT INTO sites(tracking_id, name) VALUES(?, ?)`)
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
