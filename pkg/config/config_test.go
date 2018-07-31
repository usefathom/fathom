package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	before := len(os.Environ())
	LoadEnv("")
	LoadEnv("1230")
	after := len(os.Environ())

	if before != after {
		t.Errorf("Expected the same number of env values")
	}

	data := []byte("FATHOM_DATABASE_DRIVER=\"sqlite3\"")
	ioutil.WriteFile("env_values", data, 0644)
	defer os.Remove("env_values")

	LoadEnv("env_values")

	got := os.Getenv("FATHOM_DATABASE_DRIVER")
	if got != "sqlite3" {
		t.Errorf("Expected %v, got %v", "sqlite3", got)
	}
}

func TestParse(t *testing.T) {
	// empty config, should not fatal
	cfg := Parse()
	if cfg.Secret == "" {
		t.Errorf("expected secret, got empty string")
	}

	secret := "my-super-secret-string"
	os.Setenv("FATHOM_SECRET", secret)
	cfg = Parse()
	if cfg.Secret != secret {
		t.Errorf("Expected %#v, got %#v", secret, cfg.Secret)
	}

	os.Setenv("FATHOM_DATABASE_DRIVER", "sqlite")
	cfg = Parse()
	if cfg.Database.Driver != "sqlite3" {
		t.Errorf("expected %#v, got %#v", "sqlite3", cfg.Database.Driver)
	}
}

func TestDatabaseURL(t *testing.T) {
	data := []byte("FATHOM_DATABASE_URL=\"postgres://dbuser:dbsecret@dbhost:1234/dbname\"")
	ioutil.WriteFile("env_values", data, 0644)
	defer os.Remove("env_values")

	LoadEnv("env_values")
	cfg := Parse()
	driver := "postgres"
	url := "postgres://dbuser:dbsecret@dbhost:1234/dbname"
	if cfg.Database.Driver != driver {
		t.Errorf("Expected %#v, got %#v", driver, cfg.Database.Driver)
	}
	if cfg.Database.URL != url {
		t.Errorf("Expected %#v, got %#v", url, cfg.Database.URL)
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
