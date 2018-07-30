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
