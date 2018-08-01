package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64
	Email    string
	Password string `json:"-"`
}

func NewUser(e string, pwd string) User {
	u := User{
		Email: e,
	}
	u.SetPassword(pwd)
	return u
}

func (u *User) SetPassword(pwd string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	u.Password = string(hash)
}

func (u *User) ComparePassword(pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
}
