package gorm

import (
	"database/sql"
	"fmt"

	myctx "github.com/antsrp/gdb_ex/internal/context"
	"github.com/antsrp/gdb_ex/internal/interfaces/db"
	"github.com/antsrp/gdb_ex/internal/interfaces/logger"
	repo "github.com/antsrp/gdb_ex/pkg/infrastructure/db"
	"github.com/pkg/errors"

	"gorm.io/gorm"
)

type transaction struct {
	tx *gorm.DB
}

func (t transaction) Commit() error {
	if err := t.tx.Commit().Error; err != nil {
		return fmt.Errorf("can't commit transaction: %w", err)
	}

	return nil
}

func (t transaction) Rollback() error {
	if err := t.tx.Rollback().Error; err != nil {
		return fmt.Errorf("can't rollback transaction: %w", err)
	}

	return nil
}

type connection struct {
	gdb    *gorm.DB
	logger logger.Logger
	txKey  myctx.TxKey
}

// func NewConnection(settings *repo.Settings, logger logger.Logger, txKey myctx.TxKey) (*connection, error) {
func NewConnection(settings *repo.Settings, logger logger.Logger, txKey myctx.TxKey) (*Connection, error) {
	dialect, err := withDialect(settings)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(dialect, &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to db")
	}

	c := connection{
		gdb:    db,
		logger: logger,
		txKey:  txKey,
	}

	if err := c.Check(); err != nil {
		return nil, err
	}
	c.logger.Info("db connection opened")

	//return &c, nil
	return &Connection{c}, nil
}

func (c connection) Check() error {
	if db, err := c.gdb.DB(); err != nil {
		return err
	} else {
		return db.Ping()
	}
}

func (c connection) Close() error {
	c.logger.Info("db connection closing")
	db, err := c.gdb.DB()
	if err != nil {
		c.logger.Error("can't close db connection: %v\n", err.Error())
		return err
	}
	if err := db.Close(); err != nil {
		c.logger.Error("can't close db connection: %v\n", err.Error())
		return err
	} else {
		c.logger.Info("db connection closed")
	}
	return nil
}

func (c connection) CreateTransaction() (db.Transaction, error) {
	return transaction{
		tx: c.gdb.Begin(),
	}, nil
}

func (c connection) DB() (*sql.DB, error) {
	db, err := c.gdb.DB()
	if err != nil {
		return nil, fmt.Errorf("can't get sql db connection: %w", err)
	}
	return db, nil
}

type Connection struct {
	connection
}

// var _ db.Connection = connection{}
var _ db.Connection = Connection{}
