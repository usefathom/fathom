package datastore

import (
	"database/sql"
	"github.com/dannyvankooten/ana/models"
)

var err error
var stmt *sql.Stmt
var u models.User

func GetUser(id int64) (*models.User, error) {
	stmt, err = DB.Prepare("SELECT id, email FROM users WHERE id = ? LIMIT 1")
	err = stmt.QueryRow(id).Scan(&u.ID, &u.Email)
	return &u, err
}

func GetUserByEmail(email string) (*models.User, error) {
	stmt, err = DB.Prepare("SELECT id, email, password FROM users WHERE email = ? LIMIT 1")
	err := stmt.QueryRow(email).Scan(&u.ID, &u.Email, &u.HashedPassword)
	return &u, err
}
