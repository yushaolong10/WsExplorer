package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var logger *zap.Logger
var sugar *zap.SugaredLogger

func Init(filePath string, isDebug bool) (err error) {
	cfg := zap.Config{
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:       "date",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "caller",
			MessageKey:    "msg",
			StacktraceKey: "trace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(fmt.Sprintf("%d-%d-%d %d:%d:%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
			},
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{filePath},
		ErrorOutputPaths: []string{filePath},
	}
	if isDebug {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.Development = true
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		cfg.Development = false
	}
	logger, err = cfg.Build()
	sugar = logger.Sugar()
	return
}

func Close() (err error) {
	if logger != nil {
		return logger.Sync()
	}
	return
}

func Debug(format string, v ...interface{}) {
	sugar.Debugf(format, v...)
}

func Info(format string, v ...interface{}) {
	sugar.Infof(format, v...)
}

func Warn(format string, v ...interface{}) {
	sugar.Warnf(format, v...)
}

func Error(format string, v ...interface{}) {
	sugar.Errorf(format, v...)
}
