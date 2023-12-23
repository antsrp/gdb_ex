//go:build wireinject
// +build wireinject

package injections

import (
	"database/sql"
	"io/fs"

	"github.com/antsrp/gdb_ex/internal/interfaces/db"
	"github.com/antsrp/gdb_ex/pkg/usecases/migrations/goose"
	"github.com/google/wire"
)

var (
	provideGooseMigratorSet = wire.NewSet(provideGooseMigrator, provideDBConnection)
)

func provideGooseMigrator(fs fs.FS, db *sql.DB, dialect string) (goose.MigrationTool, error) {
	return goose.NewMigrationTool(fs, db, dialect)
}

func provideDBConnection(conn db.Connection) (*sql.DB, error) {
	return conn.DB()
}

func BuildGooseMigrator(fs fs.FS, dialect string, conn db.Connection) (goose.MigrationTool, error) {
	panic(wire.Build(provideGooseMigratorSet))
}
