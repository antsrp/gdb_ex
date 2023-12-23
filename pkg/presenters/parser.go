package presenters

import (
	"context"
	"fmt"

	myctx "github.com/antsrp/gdb_ex/internal/context"
)

func ParseFromCtx[T any](ctx context.Context, key myctx.TxKey) (T, error) {
	val := ctx.Value(key)
	data, ok := val.(T)
	if !ok {
		return *new(T), fmt.Errorf("can't parse data from context or there is nothing passed")
	}

	return data, nil
}
