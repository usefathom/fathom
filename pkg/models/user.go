package models

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64
	Email    string
	Password string `json:"-"`
}

// NewUser creates a new User with the given email and password
func NewUser(e string, pwd string) User {
	u := User{
		Email: strings.ToLower(strings.TrimSpace(e)),
	}
	u.SetPassword(pwd)
	return u
}

// SetPassword sets a brcrypt encrypted password from the given plaintext pwd
func (u *User) SetPassword(pwd string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	u.Password = string(hash)
}

// ComparePassword returns true when the given plaintext password matches the encrypted pwd
func (u *User) ComparePassword(pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
}
