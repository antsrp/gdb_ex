package logger

type Logger interface {
	Info(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
	Debug(string, ...interface{})
	Panic(string, ...interface{})
	DPanic(string, ...interface{})
	Warn(string, ...interface{})
}
