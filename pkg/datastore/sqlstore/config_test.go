package sqlstore

import (
	"fmt"
	"testing"
)

func TestConfigDSN(t *testing.T) {
	c := Config{
		Driver:   "postgres",
		user:     "john",
		password: "foo",
	}
	e := fmt.Sprintf("user=%s password=%s", c.user, c.password)
	if v := c.DSN(); v != e {
		t.Errorf("Invalid DSN. Expected %s, got %s", e, v)
	}

	c = Config{
		Driver:   "postgres",
		user:     "john",
		password: "foo",
		sslmode:  "disable",
	}
	e = fmt.Sprintf("user=%s password=%s sslmode=%s", c.user, c.password, c.sslmode)
	if v := c.DSN(); v != e {
		t.Errorf("Invalid DSN. Expected %s, got %s", e, v)
	}
}

func TestConfigDbname(t *testing.T) {
	var c Config

	c = Config{
		url: "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full",
	}
	if e, v := "pqgotest", c.Dbname(); v != e {
		t.Errorf("Expected %q, got %q", e, v)
	}

	c = Config{
		url: "root@tcp(host.myhost)/mysqltest?loc=Local",
	}
	if e, v := "mysqltest", c.Dbname(); v != e {
		t.Errorf("Expected %q, got %q", e, v)
	}

	c = Config{
		url: "/mysqltest?loc=Local&parseTime=true",
	}
	if e, v := "mysqltest", c.Dbname(); v != e {
		t.Errorf("Expected %q, got %q", e, v)
	}
}
