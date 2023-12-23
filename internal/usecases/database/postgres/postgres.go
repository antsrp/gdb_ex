package postgres

import (
	"context"
	"database/sql"
	"fmt"

	myctx "github.com/antsrp/gdb_ex/internal/context"
	"github.com/antsrp/gdb_ex/internal/interfaces/db"
	"github.com/antsrp/gdb_ex/internal/interfaces/logger"
	repo "github.com/antsrp/gdb_ex/pkg/infrastructure/db"
	"github.com/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type transaction struct {
	tx pgx.Tx
}

func (t transaction) Commit() error {
	if err := t.tx.Commit(context.TODO()); err != nil {
		return fmt.Errorf("can't commit transaction: %w", err)
	}

	return nil
}

func (t transaction) Rollback() error {
	if err := t.tx.Rollback(context.TODO()); err != nil {
		return fmt.Errorf("can't rollback transaction: %w", err)
	}

	return nil
}

type connection struct {
	//pc     *pgx.Conn
	pc     *pgxpool.Pool
	logger logger.Logger
	ctx    context.Context
	txKey  myctx.TxKey
}

// func NewConnection(settings *repo.Settings, logger logger.Logger, txKey myctx.TxKey) (*connection, error) {
func NewConnection(settings *repo.Settings, logger logger.Logger, txKey myctx.TxKey) (*Connection, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		settings.User, settings.Password, settings.Host, settings.Port, settings.DBName)

	ctx := context.Background()
	/*pc, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to db")
	}*/
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "can't create connection pool")
	}

	c := connection{
		//pc:     pc,
		pc:     pool,
		logger: logger,
		ctx:    ctx,
		txKey:  txKey,
	}

	if err := c.Check(); err != nil {
		return nil, err
	}
	c.logger.Info("postgres connection opened")

	return &Connection{c}, nil
	//return &c, nil
}

func (c connection) Check() error {
	return c.pc.Ping(c.ctx)
}

func (c connection) Close() error {
	c.logger.Info("postgres connection closing")
	c.pc.Close()
	return nil
}

func (c connection) CreateTransaction() (db.Transaction, error) {
	tx, err := c.pc.Begin(c.ctx)
	if err != nil {
		return nil, fmt.Errorf("can't create transaction: %w", err)
	}

	return transaction{
		tx: tx,
	}, nil
}

func (c connection) DB() (*sql.DB, error) {
	return stdlib.OpenDBFromPool(c.pc), nil
}

type Connection struct {
	connection
}

// var _ db.Connection = connection{}
var _ db.Connection = Connection{}
