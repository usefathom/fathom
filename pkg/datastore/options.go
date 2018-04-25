package datastore

import (
	"database/sql"
)

// GetOption returns an option value by its name
func GetOption(name string) (string, error) {
	value := ""
	query := dbx.Rebind(`SELECT o.value FROM options o WHERE o.name = ? LIMIT 1`)
	err := dbx.Get(&value, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrNoResults
		}

		return "", err
	}
	return value, nil
}

// SetOption updates an option by its name
func SetOption(name string, value string) error {
	query := dbx.Rebind(`INSERT INTO options(name, value) VALUES(?, ?) ON DUPLICATE KEY UPDATE value = ?`)
	_, err := dbx.Exec(query, name, value, value)
	return err
}
