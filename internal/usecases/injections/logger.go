//go:build wireinject
// +build wireinject

package injections

import (
	"fmt"

	"github.com/antsrp/gdb_ex/internal/usecases/logger/mock"
	"github.com/antsrp/gdb_ex/internal/usecases/logger/zap"
	"github.com/google/wire"
)

var (
	msgCantCreateLogger = "can't create logger"
)

var (
	provideZapLoggerSet = wire.NewSet(
		provideZapLogger,
	)
	provideMockLoggerSet = wire.NewSet(
		provideMockLogger,
	)
)

func provideZapLogger() (*zap.Logger, error) {
	logger, err := zap.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msgCantCreateLogger, err)
	}
	return logger, nil
}

func provideMockLogger() (*mock.Logger, error) {
	logger, err := mock.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msgCantCreateLogger, err)
	}
	return logger, nil
}

func BuildZapLogger() (*zap.Logger, error) {
	panic(wire.Build(provideZapLoggerSet))
}

func BuildMockLogger() (*mock.Logger, error) {
	panic(wire.Build(provideMockLoggerSet))
}
