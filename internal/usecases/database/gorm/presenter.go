package gorm

import (
	"context"

	myctx "github.com/antsrp/gdb_ex/internal/context"
	"github.com/antsrp/gdb_ex/internal/interfaces/logger"
	"github.com/antsrp/gdb_ex/pkg/presenters"
	"gorm.io/gorm"
)

func presentOperand(ctx context.Context, logger logger.Logger, connection *gorm.DB, txKey myctx.TxKey) *gorm.DB {
	var operand *gorm.DB
	transaction, err := presenters.ParseFromCtx[transaction](ctx, txKey)
	if err != nil {
		logger.Info(err.Error())
		operand = connection
	} else {
		operand = transaction.tx
	}

	return operand
}
