package sqlstore

import (
	"strings"

	mysql "github.com/go-sql-driver/mysql"
)

type Config struct {
	URL      string `default:""`
	Driver   string `default:"sqlite3"`
	Host     string `default:""`
	User     string `default:""`
	Password string `default:""`
	Name     string `default:"fathom.db"`
	SSLMode  string `default:""`
}

func (c *Config) DSN() string {
	var dsn string

	// if FATHOM_DATABASE_URL was set, use that
	// this relies on the user to set the appropriate parameters, eg ?parseTime=true when using MySQL
	if c.URL != "" {
		return c.URL
	}

	// otherwise, generate from individual fields
	switch c.Driver {
	case POSTGRES:
		params := map[string]string{
			"host":     c.Host,
			"dbname":   c.Name,
			"user":     c.User,
			"password": c.Password,
			"sslmode":  c.SSLMode,
		}

		for k, v := range params {
			if v == "" {
				continue
			}

			dsn = dsn + k + "=" + v + " "
		}

		dsn = strings.TrimSpace(dsn)
	case MYSQL:
		mc := mysql.NewConfig()
		mc.User = c.User
		mc.Passwd = c.Password
		mc.Addr = c.Host
		mc.Net = "tcp"
		mc.DBName = c.Name
		mc.Params = map[string]string{
			"parseTime": "true",
			"loc":       "Local",
		}
		if c.SSLMode != "" {
			mc.Params["tls"] = c.SSLMode
		}
		dsn = mc.FormatDSN()
	case SQLITE:
		dsn = c.Name + "?_loc=auto&_busy_timeout=5000"
	}

	return dsn
}
