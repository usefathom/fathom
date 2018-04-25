package datastore

import (
	"github.com/usefathom/fathom/pkg/models"
)

// GetPageByHostnameAndPath retrieves a page from the connected database
func GetPageByHostnameAndPath(hostname, path string) (*models.Page, error) {
	p := &models.Page{}

	sql := dbx.Rebind(`SELECT * FROM pages WHERE hostname = ? AND path = ? LIMIT 1`)
	err := dbx.Get(p, sql, hostname, path)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// SavePage inserts the page model in the connected database
func SavePage(p *models.Page) error {
	sql := dbx.Rebind(`INSERT INTO pages(hostname, path, title) VALUES(?, ?, ?)`)
	result, err := dbx.Exec(sql, p.Hostname, p.Path, p.Title)
	if err != nil {
		return err
	}

	p.ID, _ = result.LastInsertId()
	return nil
}
