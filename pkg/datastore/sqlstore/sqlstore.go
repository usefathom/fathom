package sqlstore

import (
	"errors"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/gobuffalo/packr"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"           // postgresql driver
	_ "github.com/mattn/go-sqlite3" //sqlite3 driver
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
)

type sqlstore struct {
	*sqlx.DB

	Config *Config
}

// ErrNoResults is returned when a query yielded 0 results
var ErrNoResults = errors.New("datastore: query returned 0 results")

// New creates a new database pool
func New(c *Config) *sqlstore {
	dbx := sqlx.MustConnect(c.Driver, c.DSN())
	db := &sqlstore{dbx, c}

	// run migrations
	db.Migrate()

	return db
}

func (db *sqlstore) Migrate() {
	migrations := &migrate.PackrMigrationSource{
		Box: packr.NewBox("./migrations"),
		Dir: "./" + db.Config.Driver,
	}
	migrate.SetTable("migrations")

	n, err := migrate.Exec(db.DB.DB, db.Config.Driver, migrations, migrate.Up)

	if err != nil {
		log.Fatal("database migrations failed: ", err)
	}

	if n > 0 {
		log.Infof("applied %d database migrations", n)
	}
}

// Closes the db pool
func (db *sqlstore) Close() {
	db.DB.Close()
}
