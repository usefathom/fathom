package sqlstore

import (
	"fmt"
	"testing"
)

func TestConfigDSN(t *testing.T) {
	c := Config{
		Driver:   "postgres",
		User:     "john",
		Password: "foo",
	}
	e := fmt.Sprintf("user=%s password=%s", c.User, c.Password)
	if v := c.DSN(); v != e {
		t.Errorf("Invalid DSN. Expected %s, got %s", e, v)
	}

	c = Config{
		Driver:   "postgres",
		User:     "john",
		Password: "foo",
		SSLMode:  "disable",
	}
	e = fmt.Sprintf("user=%s password=%s sslmode=%s", c.User, c.Password, c.SSLMode)
	if v := c.DSN(); v != e {
		t.Errorf("Invalid DSN. Expected %s, got %s", e, v)
	}
}

func TestConfigDbname(t *testing.T) {
	var c Config

	c = Config{
		URL: "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full",
	}
	if e, v := "pqgotest", c.Dbname(); v != e {
		t.Errorf("Expected %q, got %q", e, v)
	}

	c = Config{
		URL: "root@tcp(host.myhost)/mysqltest?loc=Local",
	}
	if e, v := "mysqltest", c.Dbname(); v != e {
		t.Errorf("Expected %q, got %q", e, v)
	}

	c = Config{
		URL: "/mysqltest?loc=Local&parseTime=true",
	}
	if e, v := "mysqltest", c.Dbname(); v != e {
		t.Errorf("Expected %q, got %q", e, v)
	}
}
