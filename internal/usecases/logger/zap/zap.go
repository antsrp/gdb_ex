package zap

import (
	"fmt"

	ilogger "github.com/antsrp/gdb_ex/internal/interfaces/logger"
	"go.uber.org/zap"
)

var _ ilogger.Logger = logger{}

type logger struct {
	zl *zap.Logger
}

func NewLogger() (*Logger, error) {
	zl, err := zap.NewDevelopment()

	if err != nil {
		return nil, fmt.Errorf("can't init zap logger: %w", err)
	}

	return &Logger{
		logger: logger{zl: zl},
	}, nil
}

func (l logger) withArgs(args ...interface{}) bool {
	return len(args) > 0
}

func (l logger) sugared() *zap.SugaredLogger {
	return l.zl.Sugar()
}

func (l logger) Info(template string, args ...interface{}) {
	sugar := l.sugared()
	if l.withArgs(args...) {
		sugar.Infof(template, args...)
	} else {
		sugar.Info(template)
	}

}
func (l logger) Error(template string, args ...interface{}) {
	sugar := l.sugared()
	if l.withArgs(args...) {
		sugar.Errorf(template, args...)
	} else {
		sugar.Error(template)
	}
}
func (l logger) Fatal(template string, args ...interface{}) {
	defer func() {
		sugar := l.sugared()
		if l.withArgs(args...) {
			sugar.Fatalf(template, args...)
		} else {
			sugar.Fatal(template)
		}
	}()
}
func (l logger) Debug(template string, args ...interface{}) {
	sugar := l.sugared()
	if l.withArgs(args...) {
		sugar.Debugf(template, args...)
	} else {
		sugar.Debug(template)
	}
}
func (l logger) Panic(template string, args ...interface{}) {
	sugar := l.sugared()
	if l.withArgs(args...) {
		sugar.Panicf(template, args...)
	} else {
		sugar.Panic(template)
	}
}
func (l logger) DPanic(template string, args ...interface{}) {
	sugar := l.sugared()
	if l.withArgs(args...) {
		sugar.DPanicf(template, args...)
	} else {
		sugar.DPanic(template)
	}
}

func (l logger) Warn(template string, args ...interface{}) {
	sugar := l.sugared()
	if l.withArgs(args...) {
		sugar.Warnf(template, args...)
	} else {
		sugar.Warn(template)
	}
}

type Logger struct {
	logger
}
