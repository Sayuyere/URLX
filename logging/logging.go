package logging

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
	loki    *LokiClient
	batcher *LokiBatcher
}

var (
	loggerInstance *Logger
	once           sync.Once
)

func NewLogger() *Logger {
	once.Do(func() {
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		cfg.EncoderConfig.CallerKey = "caller"
		cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.AddSync(zapcore.Lock(os.Stdout)),
			cfg.Level,
		)
		zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		lokiClient := NewLokiClient()
		batcher := NewLokiBatcher(lokiClient, 10, 2*time.Second) // batch size 10, flush every 2s
		loggerInstance = &Logger{zapLogger.Sugar(), lokiClient, batcher}
	})
	return loggerInstance
}

// Helper to encode Zap log as JSON with caller info
func encodeZapJSON(level, msg string, skip int, args ...interface{}) string {
	fields := make([]zapcore.Field, 0, len(args)/2)
	for i := 0; i+1 < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		fields = append(fields, zap.Any(key, args[i+1]))
	}
	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	zapLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		zapLevel = zapcore.DebugLevel
	}
	_, file, line, ok := runtime.Caller(skip)
	buf, _ := enc.EncodeEntry(
		zapcore.Entry{
			Level:   zapLevel,
			Message: msg,
			Time:    time.Now(),
			Caller:  zapcore.EntryCaller{Defined: ok, File: file, Line: line},
		},
		fields,
	)
	return buf.String()
}

// Helper for int to string (no import strconv for one call)
func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.Infof(msg, args...)
	if l.batcher != nil {
		jsonLog := encodeZapJSON("info", msg, 2, args...)
		l.batcher.SendLog(jsonLog, "info", l.loki.serviceName, nil)
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.Errorf(msg, args...)
	if l.batcher != nil {
		jsonLog := encodeZapJSON("error", msg, 2, args...)
		l.batcher.SendLog(jsonLog, "error", l.loki.serviceName, nil)
	}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.Debugf(msg, args...)
	if l.batcher != nil {
		jsonLog := encodeZapJSON("debug", msg, 2, args...)
		l.batcher.SendLog(jsonLog, "debug", l.loki.serviceName, nil)
	}
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.Warnf(msg, args...)
	if l.batcher != nil {
		jsonLog := encodeZapJSON("warn", msg, 2, args...)
		l.batcher.SendLog(jsonLog, "warn", l.loki.serviceName, nil)
	}
}
