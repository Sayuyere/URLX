package logging

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		zapLogger, _ := cfg.Build()
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
