package logging

import (
	"sync"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

var (
	loggerInstance *Logger
	once           sync.Once
)

func NewLogger() *Logger {
	once.Do(func() {
		zapLogger, _ := zap.NewProduction()
		loggerInstance = &Logger{zapLogger.Sugar()}
	})
	return loggerInstance
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.Infof(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.Errorf(msg, args...)
}
