package soa

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type defaultLogger struct {
}

func (log *defaultLogger) Info(args ...interface{}) {}

func (log *defaultLogger) Error(args ...interface{}) {}

var logger Logger = &defaultLogger{}

func SetLogger(extraLogger Logger) {
	logger = extraLogger
}
