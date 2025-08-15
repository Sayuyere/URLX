package logging

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger() *Logger {
	logger, _ := zap.NewProduction()
	return &Logger{logger.Sugar()}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.Infof(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.Errorf(msg, args...)
}
