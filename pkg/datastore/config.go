package datastore

import "fmt"

type Config struct {
	Driver   string `default:"sqlite"`
	Host     string `default:""`
	User     string `default:""`
	Password string `default:""`
	Name     string `default:"fathom.db"`
}

func (c *Config) DSN() string {
	var dsn string

	switch c.Driver {
	case "postgres":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s", c.Host, c.User, c.Password, c.Name)
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@%s/%s?parseTime=true&loc=Local", c.User, c.Password, c.Host, c.Name)
	case "sqlite3", "sqlite":

		dsn = c.Name + "?_loc=auto" // TODO: Make this configurable
	}

	return dsn
}
