package sqlstore

import (
	"regexp"
	"strings"

	mysql "github.com/go-sql-driver/mysql"
)

type Config struct {
	Driver   string `default:"sqlite3"`
	URL      string `default:""`
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
		if c.Host != "" {
			dsn += " host=" + c.Host
		}
		if c.Name != "" {
			dsn += " dbname=" + c.Name
		}
		if c.User != "" {
			dsn += " user=" + c.User
		}
		if c.Password != "" {
			dsn += " password=" + c.Password
		}
		if c.SSLMode != "" {
			dsn += " sslmode=" + c.SSLMode
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

// Dbname returns the database name, either from config values or from the connection URL
func (c *Config) Dbname() string {
	if c.Name != "" {
		return c.Name
	}

	re := regexp.MustCompile(`(?:dbname=|[^\/]?\/)(\w+)`)
	m := re.FindStringSubmatch(c.URL)
	if len(m) > 1 {
		return m[1]
	}

	return ""
}
