package datastore

import (
	"github.com/dannyvankooten/ana/models"
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

	err = stmt.QueryRow(email).Scan(&u.ID, &u.Email, &u.HashedPassword)
	return &u, err
}
