package migrator

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *Migrator) Up() error {
	// goose.Up применяет все миграции, которые еще не были применены
	err := goose.Up(m.db, m.migrationsDir)
	if err != nil {
		return err
	}

	return nil
}
