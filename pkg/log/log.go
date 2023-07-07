package log

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var loggerKey struct{}

var logger *zap.Logger

type Config struct {
	MinLevel zapcore.Level
	Stdout   bool

	// Filename is the file to write logs to. Backup log files will be retained
	// in the same directory.
	Filename string
	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int
	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int
}

func Init(cfg *Config) *zap.Logger {
	if cfg.MaxSize == 0 {
		cfg.MaxSize = 100 // 100 mb
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 7 // 7 days
	}
	core := newZapCore(cfg)
	logger = zap.New(
		core,
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	return logger
}

func newZapCore(cfg *Config) zapcore.Core {
	opts := []zapcore.WriteSyncer{
		zapcore.AddSync(&lumberjack.Logger{
			Filename: cfg.Filename,
			MaxSize:  cfg.MaxSize,
			MaxAge:   cfg.MaxAge,
		}),
	}
	if cfg.Stdout {
		opts = append(opts, zapcore.AddSync(zapcore.Lock(os.Stdout)))
	}
	syncWriter := zapcore.NewMultiWriteSyncer(opts...)

	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(level.CapitalString())
	}
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(caller.TrimmedPath())
	}
	encoder := zapcore.EncoderConfig{
		CallerKey:      "line",
		LevelKey:       "level_name",
		MessageKey:     "api",
		TimeKey:        "log_time",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     customTimeEncoder,   // 自定义时间格式
		EncodeLevel:    customLevelEncoder,  // 小写编码器
		EncodeCaller:   customCallerEncoder, // 全路径编码器
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewCore(zapcore.NewJSONEncoder(encoder), syncWriter, zap.NewAtomicLevelAt(cfg.MinLevel))
}

func WithCtx(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(loggerKey).(*zap.Logger)
	if ok {
		return l.WithOptions(zap.AddCallerSkip(-1))
	}
	return logger
}

func AddCtx(ctx context.Context, fields ...zap.Field) context.Context {
	return context.WithValue(ctx, loggerKey, logger.With(fields...))
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
