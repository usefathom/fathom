package api

import "testing"

func TestLoginSanitize(t *testing.T) {
	rawEmail := "Foo@foobar.com   "
	l := &login{
		Email: rawEmail,
	}

	l.Sanitize()
	if l.Email != "foo@foobar.com" {
		t.Errorf("Expected normalized email address, got %s", l.Email)
	}
}
