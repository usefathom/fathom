package datastore

// GetOption returns an option value by its name
func GetOption(name string) (string, error) {
	var value string

	stmt, err := DB.Prepare(`SELECT o.value FROM options o WHERE o.name = ? LIMIT 1`)
	if err != nil {
		return "", err
	}

	defer stmt.Close()
	stmt.QueryRow(name).Scan(&value)

	return value, nil
}

// SetOption updates an option by its name
func SetOption(name string, value string) error {
	stmt, err := DB.Prepare(`INSERT INTO options(name, value) VALUES(?, ?) ON DUPLICATE KEY UPDATE value = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(name, value, value)

	return err
}
