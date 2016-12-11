package models

import (
	"database/sql"
)

type User struct {
	ID       int64
	Email    string
	Password string `json:"-"`
}

func (u *User) Save(conn *sql.DB) error {
	// prepare statement for inserting data
	stmt, err := conn.Prepare(`INSERT INTO users(
    email,
    password
    ) VALUES(?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.Email, u.Password)
	u.ID, _ = result.LastInsertId()

	return err
}
