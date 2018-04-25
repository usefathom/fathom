package datastore

// GetOption returns an option value by its name
func GetOption(name string) (string, error) {
	value := ""
	err := dbx.Get(&value, dbx.Rebind(`SELECT o.value FROM options o WHERE o.name = ? LIMIT 1`), name)
	if err != nil {
		return "", err
	}
	return value, nil
}

// SetOption updates an option by its name
func SetOption(name string, value string) error {
	sql := dbx.Rebind(`INSERT INTO options(name, value) VALUES(?, ?) ON DUPLICATE KEY UPDATE value = ?`)
	_, err := dbx.Exec(sql, name, value, value)
	return err
}
