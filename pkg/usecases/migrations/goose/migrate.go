package goose

import (
	"database/sql"
	"fmt"
	"io/fs"

	imig "github.com/antsrp/gdb_ex/internal/interfaces/migrations"
	"github.com/pressly/goose/v3"
)

var _ imig.Migrator = MigrationTool{}

type migrationTool struct {
	db *sql.DB
}

func NewMigrationTool(fs fs.FS, db *sql.DB, dialect string) (MigrationTool, error) {
	goose.SetBaseFS(fs)
	if err := goose.SetDialect(dialect); err != nil {
		return MigrationTool{}, fmt.Errorf("can't create migration tool: %w", err)
	}

	return MigrationTool{
		migrationTool{db: db},
	}, nil
}

func (m migrationTool) Up(folder string) error {
	if err := goose.Up(m.db, folder); err != nil {
		return fmt.Errorf("can't up migration: %w", err)
	}
	return nil
}

func (m migrationTool) Down(folder string) error {
	if err := goose.Down(m.db, folder); err != nil {
		return fmt.Errorf("can't down migration: %w", err)
	}
	return nil
}

func (m migrationTool) UpTo(folder string, version int64) error {
	if err := goose.UpTo(m.db, folder, version); err != nil {
		return fmt.Errorf("can't up to migration %d: %w", version, err)
	}
	return nil
}

func (m migrationTool) DownTo(folder string, version int64) error {
	if err := goose.DownTo(m.db, folder, version); err != nil {
		return fmt.Errorf("can't down to migration %d: %w", version, err)
	}
	return nil
}

func (m migrationTool) DownAll(folder string) error {
	return m.DownTo(folder, 0)
}

type MigrationTool struct {
	migrationTool
}
