package log

import (
	"context"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *otelzap.Logger

func init() {
	l, err := zap.NewProductionConfig().Build(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		panic(err)
	}
	logger = otelzap.New(l, otelzap.WithTraceIDField(true))
}

func WithContext(ctx context.Context) otelzap.LoggerWithCtx {
	return logger.Ctx(ctx)
}

func WithOptions(opts ...zap.Option) *zap.Logger {
	return logger.Logger.WithOptions(opts...)
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
