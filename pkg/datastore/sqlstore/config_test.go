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
