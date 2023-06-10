package mlog

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var appName = "minerva"

var appLogger = log.With(Logger(appName),
	"ts", log.DefaultTimestamp,
	"caller", log.DefaultCaller,
	"service.id", appName,
	"service", appName,
	"service.version", "0.1",
	"trace.id", tracing.TraceID(),
	"span.id", tracing.SpanID(),
)

var StdLogHelper = log.NewHelper(appLogger)

type ZapLogger struct {
	log  *zap.Logger
	Sync func() error
}

func Logger(appName string) log.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
	return NewZLogger(
		appName,
		encoderConfig,
		zap.NewAtomicLevelAt(zapcore.DebugLevel),
		zap.AddStacktrace(zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.Development(),
	)
}

func logWriteSyncer(appName string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("/data/logs/%s/debug.log", appName),
		MaxSize:    1000,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func NewZLogger(appName string, encoder zapcore.EncoderConfig,
	level zap.AtomicLevel, opts ...zap.Option) *ZapLogger {
	writeSyncer := logWriteSyncer(appName)
	level.SetLevel(zap.InfoLevel)
	var core zapcore.Core
	core = zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(writeSyncer)),
		level,
	)
	zLogger := zap.New(core, opts...)
	return &ZapLogger{log: zLogger, Sync: zLogger.Sync}
}

func (l *ZapLogger) Log(level log.Level, kvs ...interface{}) error {
	if len(kvs) == 0 || len(kvs)%2 != 0 {
		l.log.Warn(fmt.Sprint("Key values must appear in pairs: ", kvs))
		return nil
	}

	var fields []zap.Field
	for i := 0; i < len(kvs); i += 2 {
		fields = append(fields, zap.Any(fmt.Sprint(kvs[i]), kvs[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug("", fields...)
	case log.LevelInfo:
		l.log.Info("", fields...)
	case log.LevelWarn:
		l.log.Warn("", fields...)
	case log.LevelError:
		l.log.Error("", fields...)
	case log.LevelFatal:
		l.log.Fatal("", fields...)
	}
	return nil
}
