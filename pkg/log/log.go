package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewDevelopment(zap.AddStacktrace(zapcore.FatalLevel), zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func WithOptions(opts ...zap.Option) *zap.Logger {
	return logger.WithOptions(opts...)
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Sync() {
	logger.Sync()
}
