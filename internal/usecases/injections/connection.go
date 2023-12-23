//go:build wireinject
// +build wireinject

package injections

import (
	"fmt"

	myctx "github.com/antsrp/gdb_ex/internal/context"
	"github.com/antsrp/gdb_ex/internal/interfaces/db"
	"github.com/antsrp/gdb_ex/internal/interfaces/logger"
	dbs "github.com/antsrp/gdb_ex/internal/usecases/database"
	"github.com/antsrp/gdb_ex/internal/usecases/database/gorm"
	"github.com/antsrp/gdb_ex/internal/usecases/database/postgres"
	"github.com/antsrp/gdb_ex/internal/usecases/service"
	idb "github.com/antsrp/gdb_ex/pkg/infrastructure/db"
	_ "github.com/google/subcommands"
	"github.com/google/wire"
)

var (
	msgCantInitConnection = "can't init database connection"
	msgCantPingConnection = "can't ping database connection"

	msgCantInitSettings = "can't init database connection settings"
)

var (
	provideGormConnectionSet = wire.NewSet(
		provideGormConnection,
	)
	providePgxConnectionSet = wire.NewSet(
		providePgxConnection,
	)
	providePgxSet = wire.NewSet(
		service.NewService,
		postgres.NewDepartmentStorage,
		postgres.NewEmployeeStorage,
		postgres.NewProjectStorage,
		wire.Bind(new(db.DepartmentsRepository), new(*postgres.DepartmentStorage)),
		wire.Bind(new(db.EmployeesRepository), new(*postgres.EmployeeStorage)),
		wire.Bind(new(db.ProjectsRepository), new(*postgres.ProjectStorage)),
	)
	provideGormSet = wire.NewSet(
		service.NewService,
		gorm.NewDepartmentStorage,
		gorm.NewEmployeeStorage,
		gorm.NewProjectStorage,
		wire.Bind(new(db.DepartmentsRepository), new(*gorm.DepartmentStorage)),
		wire.Bind(new(db.EmployeesRepository), new(*gorm.EmployeeStorage)),
		wire.Bind(new(db.ProjectsRepository), new(*gorm.ProjectStorage)),
	)
	provideSettingsSet = wire.NewSet(
		provideSettings,
	)
	provideGormStoragesSet = wire.NewSet(
		provideGormStorages,
		gorm.NewDepartmentStorage,
		gorm.NewEmployeeStorage,
		gorm.NewProjectStorage,
	)
	providePgxStoragesSet = wire.NewSet(
		providePgxStorages,
		postgres.NewDepartmentStorage,
		postgres.NewEmployeeStorage,
		postgres.NewProjectStorage,
	)
)

type PgxStorages struct {
	Ds *postgres.DepartmentStorage
	Es *postgres.EmployeeStorage
	Ps *postgres.ProjectStorage
}

type GormStorages struct {
	Ds *gorm.DepartmentStorage
	Es *gorm.EmployeeStorage
	Ps *gorm.ProjectStorage
}

func providePgxStorages(ds *postgres.DepartmentStorage, es *postgres.EmployeeStorage, ps *postgres.ProjectStorage) PgxStorages {
	return PgxStorages{
		Ds: ds,
		Es: es,
		Ps: ps,
	}
}

func provideGormStorages(ds *gorm.DepartmentStorage, es *gorm.EmployeeStorage, ps *gorm.ProjectStorage) GormStorages {
	return GormStorages{
		Ds: ds,
		Es: es,
		Ps: ps,
	}
}

func provideGormConnection(settings *idb.Settings, lg logger.Logger, txKey myctx.TxKey) (*gorm.Connection, error) {
	cn, err := gorm.NewConnection(settings, lg, txKey)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msgCantInitConnection, err)
	}
	if err = cn.Check(); err != nil {
		return nil, fmt.Errorf("%s: %w", msgCantPingConnection, err)
	}
	return cn, nil
}

func providePgxConnection(settings *idb.Settings, lg logger.Logger, txKey myctx.TxKey) (*postgres.Connection, error) {
	cn, err := postgres.NewConnection(settings, lg, txKey)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msgCantInitConnection, err)
	}
	if err = cn.Check(); err != nil {
		return nil, fmt.Errorf("%s: %w", msgCantPingConnection, err)
	}
	return cn, nil
}

func provideSettings(prefix string, filenames ...string) (*idb.Settings, error) {
	dbSettings, err := dbs.InitSettings(prefix, filenames...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msgCantInitSettings, err)
	}
	return dbSettings, nil
}

func BuildServicePgx(conn *postgres.Connection, lg logger.Logger, txKey myctx.TxKey) (*service.Service, func(), error) {
	panic(wire.Build(providePgxSet, wire.Bind(new(db.Connection), new(*postgres.Connection))))
}

func BuildServiceGorm(conn *gorm.Connection, lg logger.Logger, txKey myctx.TxKey) (*service.Service, func(), error) {
	panic(wire.Build(provideGormSet, wire.Bind(new(db.Connection), new(*gorm.Connection))))
}

func BuildConnectionPgx(settings *idb.Settings, lg logger.Logger, txKey myctx.TxKey) (*postgres.Connection, error) {
	panic(wire.Build(providePgxConnectionSet))
}

func BuildConnectionGorm(settings *idb.Settings, lg logger.Logger, txKey myctx.TxKey) (*gorm.Connection, error) {
	panic(wire.Build(provideGormConnectionSet))
}

func BuildConnectionSettings(prefix string, filenames ...string) (*idb.Settings, error) {
	panic(wire.Build(provideSettingsSet))
}

func BuildStoragesPgx(conn *postgres.Connection) (PgxStorages, error) {
	panic(wire.Build(providePgxStoragesSet))
}

func BuildStoragesGorm(conn *gorm.Connection) (GormStorages, error) {
	panic(wire.Build(provideGormStoragesSet))
}
