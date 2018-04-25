package datastore

import (
	"database/sql"
	"github.com/usefathom/fathom/pkg/models"
)

// GetUser retrieves user from datastore by its ID
func GetUser(ID int64) (*models.User, error) {
	u := &models.User{}
	query := dbx.Rebind("SELECT * FROM users WHERE id = ? LIMIT 1")
	err := dbx.Get(u, query, ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResults
		}

		return nil, err
	}

	return u, err
}

// GetUserByEmail retrieves user from datastore by its email
func GetUserByEmail(email string) (*models.User, error) {
	u := &models.User{}
	query := dbx.Rebind("SELECT * FROM users WHERE email = ? LIMIT 1")
	err := dbx.Get(u, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResults
		}

		return nil, err
	}

	return u, err
}

// SaveUser inserts the user model in the connected database
func SaveUser(u *models.User) error {
	var query = dbx.Rebind("INSERT INTO users(email, password) VALUES(?, ?)")
	result, err := dbx.Exec(query, u.Email, u.Password)
	if err != nil {
		return err
	}

	u.ID, _ = result.LastInsertId()
	return nil
}
