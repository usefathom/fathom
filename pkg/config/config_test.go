package config

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	// empty config, should not fatal
	cfg := Parse("")
	if cfg.Secret == "" {
		t.Errorf("expected secret, got empty string")
	}

	os.Setenv("FATHOM_DATABASE_DRIVER", "sqlite")
	cfg = Parse("")
	if cfg.Database.Driver != "sqlite3" {
		t.Errorf("expected %#v, got %#v", "sqlite3", cfg.Database.Driver)
	}

}

func TestRandomString(t *testing.T) {
	r1 := randomString(10)
	r2 := randomString(10)

	if r1 == r2 {
		t.Errorf("expected two different strings, got %#v", r1)
	}

	if l := len(r1); l != 10 {
		t.Errorf("expected string of length %d, got string of length %d", 10, l)
	}
}
