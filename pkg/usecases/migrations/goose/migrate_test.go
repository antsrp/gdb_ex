package goose_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	source "github.com/antsrp/gdb_ex"
	myctx "github.com/antsrp/gdb_ex/internal/context"
	"github.com/antsrp/gdb_ex/internal/interfaces/migrations"
	"github.com/antsrp/gdb_ex/internal/usecases/injections"
	"github.com/antsrp/gdb_ex/pkg/helpers"
)

var (
	pgxMigrationTool, gormMigrationTool migrations.Migrator
	folder                              string = "migrations/postgres"
)

func TestMain(m *testing.M) {
	logger, err := injections.BuildMockLogger()
	if err != nil {
		log.Fatal(err)
	}
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("path: %v\n", path)

	settings, err := injections.BuildConnectionSettings("DB", "test.env")
	if err != nil {
		logger.Fatal(err.Error())
	}

	var txKey myctx.TxKey = "tx"
	pgxConnection, err := injections.BuildConnectionPgx(settings, logger, txKey)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer helpers.HandleCloser(logger, "database pgx connection", pgxConnection)

	pgxMigrationTool, err = injections.BuildGooseMigrator(source.Migrations, settings.Type, pgxConnection)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	gormConnection, err := injections.BuildConnectionPgx(settings, logger, txKey)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer helpers.HandleCloser(logger, "database gorm connection", gormConnection)

	gormMigrationTool, err = injections.BuildGooseMigrator(source.Migrations, settings.Type, gormConnection)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	m.Run()
}
