package sqlstore

import (
	"database/sql"
)

func (db *sqlstore) HostnameID(name string) (int64, error) {
	var id int64
	query := db.Rebind("SELECT id FROM hostnames WHERE name = ? LIMIT 1")
	err := db.Get(&id, query, name)

	if err == sql.ErrNoRows {
		// Postgres does not support LastInsertID, so use a "... RETURNING" select query
		query := db.Rebind(`INSERT INTO hostnames(name) VALUES(?)`)
		if db.Driver == POSTGRES {
			err := db.Get(&id, query+" RETURNING id", name)
			return id, err
		}

		// MySQL and SQLite do support LastInsertID, so use that
		r, err := db.Exec(query, name)
		if err != nil {
			return 0, err
		}

		return r.LastInsertId()
	}

	return id, err
}
