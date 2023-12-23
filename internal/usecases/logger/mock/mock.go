package mock

import (
	ilogger "github.com/antsrp/gdb_ex/internal/interfaces/logger"
)

var _ ilogger.Logger = logger{}

type logger struct {
}

func NewLogger() (*Logger, error) {
	return &Logger{logger{}}, nil
}

func (l logger) Info(template string, args ...interface{}) {
}
func (l logger) Error(template string, args ...interface{}) {
}
func (l logger) Fatal(template string, args ...interface{}) {
}
func (l logger) Debug(template string, args ...interface{}) {
}
func (l logger) Panic(template string, args ...interface{}) {
}
func (l logger) DPanic(template string, args ...interface{}) {
}

func (l logger) Warn(template string, args ...interface{}) {
}

type Logger struct {
	logger
}
