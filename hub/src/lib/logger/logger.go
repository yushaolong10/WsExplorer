package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

func Init(filePath string, isDebug bool) error {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filePath, // ⽇志⽂件路径
		MaxSize:    4096,     // megabytes
		MaxBackups: 3,        // 最多保留3个备份
		MaxAge:     7,        //days
		Compress:   true,     // 是否压缩 disabled by default
	})
	var level zapcore.Level
	if isDebug {
		level = zap.DebugLevel
	} else {
		level = zap.InfoLevel
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		level,
	)
	logger := zap.New(core)
	sugar = logger.Sugar()
	return nil
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
