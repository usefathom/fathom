package sqlstore

import (
	mysql "github.com/go-sql-driver/mysql"
	"strings"
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

	switch c.Driver {
	case "postgres":
		if c.URL != "" {
			dsn = c.URL
		} else {
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
		}
	case "mysql":
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
	case "sqlite3":
		dsn = c.Name + "?_loc=auto"
	}

	return dsn
}
