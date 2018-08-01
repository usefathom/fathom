package models

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	email := "foo@bar.com"
	pwd := "passw0rd01"
	u := NewUser(email, pwd)

	if u.Email != email {
		t.Errorf("Email: expected %s, got %s", email, u.Email)
	}

	if u.ComparePassword(pwd) != nil {
		t.Error("Password not set correctly")
	}
}

func TestUserPassword(t *testing.T) {
	u := &User{}
	u.SetPassword("password")
	if u.ComparePassword("password") != nil {
		t.Errorf("Password should match, but does not")
	}
}
