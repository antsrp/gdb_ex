package db

import "database/sql"

type Transaction interface {
	Rollback() error
	Commit() error
}

type Connection interface {
	Check() error
	Close() error
	CreateTransaction() (Transaction, error)
	DB() (*sql.DB, error)
}

type OperationOptions struct {
	Tx Transaction
}
