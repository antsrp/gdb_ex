package gorm

import (
	"errors"
	"fmt"

	repo "github.com/antsrp/gdb_ex/pkg/infrastructure/db"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	ErrorWrongDialect = errors.New("wrong dialect")
)

func withDialect(settings *repo.Settings) (gorm.Dialector, error) {
	switch settings.Type {
	case repo.TypeMySql:
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)//%s",
			settings.User, settings.Password, settings.Host, settings.Port, settings.DBName)
		return mysql.Open(dsn), nil
	case repo.TypePostgres:
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			settings.User, settings.Password, settings.Host, settings.Port, settings.DBName)
		return postgres.Open(dsn), nil
	case repo.TypeSqlite:
		dsn := "gorm.db" // to do
		return sqlite.Open(dsn), nil
	case repo.TypeSqlServer:
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
			settings.User, settings.Password, settings.Host, settings.Port, settings.DBName)
		return sqlserver.Open(dsn), nil
	default:
		return nil, ErrorWrongDialect
	}
}
