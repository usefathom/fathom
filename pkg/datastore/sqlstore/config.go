package sqlstore

import (
	"fmt"
	mysql "github.com/go-sql-driver/mysql"
)

type Config struct {
	Driver   string `default:"sqlite3"`
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
		dsn = mc.FormatDSN()
	case "sqlite3":
		dsn = c.Name + "?_loc=auto"
	}

	return dsn
}
