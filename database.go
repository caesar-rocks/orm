package orm

import (
	"context"
	"database/sql"

	"github.com/charmbracelet/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
)

type DBMS string

const (
	MySQL    DBMS = "mysql"
	Postgres DBMS = "postgres"
	SQLite   DBMS = "sqlite"
)

type DatabaseConfig struct {
	DBMS  DBMS
	DSN   string
	Debug bool
}

type Database struct {
	*bun.DB
}

// NewDatabase creates a new bun.DB instance.
func NewDatabase(cfg *DatabaseConfig) *Database {
	var db *bun.DB

	switch cfg.DBMS {
	case MySQL:
		sqldb, err := sql.Open("mysql", cfg.DSN)
		if err != nil {
			log.Fatal(err)
		}
		db = bun.NewDB(sqldb, mysqldialect.New())
	case Postgres:
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DSN)))
		db = bun.NewDB(sqldb, pgdialect.New())
	case SQLite:
		sqldb, err := sql.Open(sqliteshim.ShimName, cfg.DSN)
		if err != nil {
			log.Fatal(err)
		}
		db = bun.NewDB(sqldb, sqlitedialect.New())
	}

	if cfg.Debug {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}

	return &Database{db}
}

func (db *Database) Migrate(migrations *migrate.Migrations) {
	ctx := context.Background()

	migrator := migrate.NewMigrator(db.DB, migrations)
	migrator.Init(ctx)

	if _, err := migrator.Migrate(ctx); err != nil {
		log.Fatal(err)
	}
	log.Info("Migrations applied successfully")
}

func (db *Database) Reset(migrations *migrate.Migrations) {
	ctx := context.Background()

	migrator := migrate.NewMigrator(db.DB, migrations)
	migrator.Init(ctx)

	if err := migrator.Reset(ctx); err != nil {
		log.Fatal(err)
	}
	log.Info("Migrations reset successfully")
}

func (db *Database) Rollback(migrations *migrate.Migrations) {
	ctx := context.Background()

	migrator := migrate.NewMigrator(db.DB, migrations)
	migrator.Init(ctx)

	if _, err := migrator.Rollback(ctx); err != nil {
		log.Fatal(err)
	}
	log.Info("Migrations rolled back successfully")
}

func (db *Database) Seed() {}
