package datastore

import (
	"github.com/dannyvankooten/ana/pkg/models"
)

var p models.Page

// GetPage ...
func GetPage(id int64) (*models.Page, error) {
	return &p, nil
}

// GetPageByHostnameAndPath ...
func GetPageByHostnameAndPath(hostname, path string) (*models.Page, error) {
	stmt, err := DB.Prepare("SELECT p.id, p.hostname, p.path FROM pages p WHERE p.hostname = ? AND p.path = ? LIMIT 1")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	err = stmt.QueryRow(hostname, path).Scan(&p.ID, &p.Hostname, &p.Path)
	return &p, err
}

// SavePage ...
func SavePage(p *models.Page) error {
	// prepare statement for inserting data
	stmt, err := DB.Prepare(`INSERT INTO pages(
			hostname,
			path,
			title
			) VALUES( ?, ?, ? )`)
	defer stmt.Close()
	if err != nil {
		return err
	}

	result, err := stmt.Exec(p.Hostname, p.Path, p.Title)
	if err != nil {
		return err
	}

	p.ID, err = result.LastInsertId()
	return err
}
