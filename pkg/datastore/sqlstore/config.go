package sqlstore

import (
	"strings"

	mysql "github.com/go-sql-driver/mysql"
)

type Config struct {
	Driver   string `default:"sqlite3"`
	url      string `default:""`
	host     string `default:""`
	user     string `default:""`
	password string `default:""`
	name     string `default:"fathom.db"`
	sslmode  string `default:""`
}

func (c *Config) DSN() string {
	var dsn string

	// if FATHOM_DATABASE_URL was set, use that
	// this relies on the user to set the appropriate parameters, eg ?parseTime=true when using MySQL
	if c.url != "" {
		return c.url
	}

	// otherwise, generate from individual fields
	switch c.Driver {
	case POSTGRES:

		if c.host != "" {
			dsn += " host=" + c.host
		}
		if c.name != "" {
			dsn += " dbname=" + c.name
		}
		if c.user != "" {
			dsn += " user=" + c.user
		}
		if c.password != "" {
			dsn += " password=" + c.password
		}
		if c.sslmode != "" {
			dsn += " sslmode=" + c.sslmode
		}

		dsn = strings.TrimSpace(dsn)
	case MYSQL:
		mc := mysql.NewConfig()
		mc.User = c.user
		mc.Passwd = c.password
		mc.Addr = c.host
		mc.Net = "tcp"
		mc.DBName = c.name
		mc.Params = map[string]string{
			"parseTime": "true",
			"loc":       "Local",
		}
		if c.sslmode != "" {
			mc.Params["tls"] = c.sslmode
		}
		dsn = mc.FormatDSN()
	case SQLITE:
		dsn = c.name + "?_loc=auto&_busy_timeout=5000"
	}

	return dsn
}
