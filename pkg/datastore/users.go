package datastore

import (
	"github.com/usefathom/fathom/pkg/models"
)

var u models.User

// GetUser retrieves user from datastore by its ID
func GetUser(id int64) (*models.User, error) {
	stmt, err := DB.Prepare("SELECT id, email FROM users WHERE id = ? LIMIT 1")
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRow(id).Scan(&u.ID, &u.Email)
	return &u, err
}

// GetUserByEmail retrieves user from datastore by its email
func GetUserByEmail(email string) (*models.User, error) {
	stmt, err := DB.Prepare("SELECT id, email, password FROM users WHERE email = ? LIMIT 1")
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRow(email).Scan(&u.ID, &u.Email, &u.Password)
	return &u, err
}

// SaveUser inserts the user model in the connected database
func SaveUser(u *models.User) error {
	var sql = dbx.Rebind("INSERT INTO users(email, password) VALUES(?, ?)")
	result, err := dbx.Exec(sql, u.Email, u.Password)
	if err != nil {
		return err
	}

	u.ID, _ = result.LastInsertId()
	return nil
}
