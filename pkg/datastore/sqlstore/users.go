package sqlstore

import (
	"database/sql"

	"github.com/usefathom/fathom/pkg/models"
)

// GetUser retrieves user from datastore by its ID
func (db *sqlstore) GetUser(ID int64) (*models.User, error) {
	u := &models.User{}
	query := db.Rebind("SELECT * FROM users WHERE id = ? LIMIT 1")
	err := db.Get(u, query, ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResults
		}

		return nil, err
	}

	return u, err
}

// GetUserByEmail retrieves user from datastore by its email
func (db *sqlstore) GetUserByEmail(email string) (*models.User, error) {
	u := &models.User{}
	query := db.Rebind("SELECT * FROM users WHERE email = ? LIMIT 1")
	err := db.Get(u, query, email)
	return u, mapError(err)
}

// SaveUser inserts the user model in the connected database
func (db *sqlstore) SaveUser(u *models.User) error {
	if u.ID > 0 {
		return db.updateUser(u)
	}

	return db.insertUser(u)
}

// insertUser saves a new user in the database
func (db *sqlstore) insertUser(u *models.User) error {
	var query = db.Rebind("INSERT INTO users(email, password) VALUES(?, ?)")

	// Postgres does not support LastInsertID, so use a "... RETURNING" select query
	if db.Driver == POSTGRES {
		err := db.Get(&u.ID, query+" RETURNING id", u.Email, u.Password)
		return err
	}

	// MySQL and SQLite don't support RETURNING, but do support LastInsertId
	result, err := db.Exec(query, u.Email, u.Password)
	if err != nil {
		return err
	}

	u.ID, err = result.LastInsertId()
	return err
}

// updateUser updates an existing user in the database
func (db *sqlstore) updateUser(u *models.User) error {
	var query = db.Rebind("UPDATE users SET email = ?, password = ? WHERE id = ?")
	_, err := db.Exec(query, u.Email, u.Password, u.ID)
	return err
}

// DeleteUser deletes the user in the datastore
func (db *sqlstore) DeleteUser(u *models.User) error {
	query := db.Rebind("DELETE FROM users WHERE id = ?")
	_, err := db.Exec(query, u.ID)
	return err
}

// CountUsers returns the number of users
func (db *sqlstore) CountUsers() (int64, error) {
	var c int64
	var sql = `SELECT COUNT(*) FROM users`
	query := db.Rebind(sql)
	err := db.Get(&c, query)
	return c, err
}
