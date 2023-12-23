package postgres

import (
	"context"

	myctx "github.com/antsrp/gdb_ex/internal/context"
	"github.com/antsrp/gdb_ex/internal/interfaces/logger"
	"github.com/antsrp/gdb_ex/pkg/presenters"
	"github.com/jackc/pgx/v5/pgxpool"
)

func presentOperand(ctx context.Context, logger logger.Logger, connection *pgxpool.Pool, txKey myctx.TxKey) operand {
	var operand operand
	ts, err := presenters.ParseFromCtx[transaction](ctx, txKey)
	if err != nil {
		logger.Info(err.Error())
		operand = connection
	} else {
		operand = ts.tx
	}

	return operand
}
