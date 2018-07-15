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

	// write log statement
	log.Printf("Connected to %s database: %s", c.Driver, c.Name)

	// run migrations
	db.Migrate()

	return db
}

func (db *sqlstore) Migrate() {
	migrationSource := &migrate.PackrMigrationSource{
		Box: packr.NewBox("./migrations"),
		Dir: db.Config.Driver,
	}
	migrate.SetTable("migrations")

	migrations, err := migrationSource.FindMigrations()
	if err != nil {
		log.Errorf("Error loading database migrations: %s", err)
	}

	if len(migrations) == 0 {
		log.Fatalf("Missing database migrations")
	}

	n, err := migrate.Exec(db.DB.DB, db.Config.Driver, migrationSource, migrate.Up)
	if err != nil {
		log.Errorf("Error applying database migrations: %s", err)
	}

	if n > 0 {
		log.Infof("Applied %d database migrations!", n)
	}
}

// Closes the db pool
func (db *sqlstore) Close() error {
	return db.DB.Close()
}
