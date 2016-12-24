package options

import (
	"github.com/dannyvankooten/ana/db"
)

// Get returns an option value by its name
func Get(name string) string {
	var value string

	stmt, _ := db.Conn.Prepare(`SELECT o.value FROM options o WHERE o.name = ?`)
	defer stmt.Close()
	stmt.QueryRow(name).Scan(&value)

	return value
}

// Set updates an option by its name
func Set(name string, value string) error {
	stmt, err := db.Conn.Prepare(`INSERT INTO options(name, value) VALUES(?, ?) ON DUPLICATE KEY UPDATE value = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(name, value, value)

	return err
}
