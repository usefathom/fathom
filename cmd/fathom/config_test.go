package main

import (
	"os"
	"testing"
)

func TestParseConfig(t *testing.T) {
	// empty config, should not fatal
	cfg := parseConfig("")
	if cfg.Secret == "" {
		t.Errorf("expected secret, got empty string")
	}

	os.Setenv("FATHOM_DATABASE_DRIVER", "sqlite")
	cfg = parseConfig("")
	if cfg.Database.Driver != "sqlite3" {
		t.Errorf("expected %#v, got %#v", "sqlite3", cfg.Database.Driver)
	}

}
