package orm

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

type MigrationsStore struct {
	Migrations *migrate.Migrations
}

func InitMigrationsStore() *MigrationsStore {
	migrations := migrate.NewMigrations()

	if err := migrations.DiscoverCaller(); err != nil {
		log.Fatal(err)
	}

	return &MigrationsStore{Migrations: migrations}
}

type Migration struct {
	MigrationInterface

	Ctx context.Context
	DB  *bun.DB
}

type MigrationInterface interface {
	Up() error
	Down() error
}

func (store *MigrationsStore) RegisterMigration(m *Migration) {
	store.Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		m.Ctx = ctx
		m.DB = db
		return m.Up()
	}, func(ctx context.Context, db *bun.DB) error {
		m.Ctx = ctx
		m.DB = db
		return m.Down()
	})
}
