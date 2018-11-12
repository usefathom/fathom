package sqlstore

import (
	"database/sql"
)

func (db *sqlstore) PathnameID(name string) (int64, error) {
	var id int64
	query := db.Rebind("SELECT id FROM pathnames WHERE name = ? LIMIT 1")
	err := db.Get(&id, query, name)

	if err == sql.ErrNoRows {
		// Postgres does not support LastInsertID, so use a "... RETURNING" select query
		query := db.Rebind(`INSERT INTO pathnames(name) VALUES(?)`)
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
