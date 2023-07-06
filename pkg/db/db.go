package db

import (
	"context"
	"errors"
	"fmt"
	"master-otel/pkg/log"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

func (c *Config) validate() error {
	if c.Host == "" {
		return errors.New("db host is required")
	}
	if c.Port == "" {
		return errors.New("db port is required")
	}
	if c.User == "" {
		return errors.New("db user is required")
	}
	if c.Password == "" {
		return errors.New("db password is required")
	}
	if c.DB == "" {
		return errors.New("db name is required")
	}
	return nil
}

func New(cfg *Config) (*gorm.DB, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DB,
		cfg.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: &customLogger{},
	})
	if err != nil {
		return nil, fmt.Errorf("gorm open: %w", err)
	}
	raw, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("raw db: %w", err)
	}
	if err := raw.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return db, err
}

type customLogger struct {
}

// LogMode log mode
func (l *customLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info print info
func (l customLogger) Info(ctx context.Context, msg string, data ...interface{}) {
}

// Warn print warn messages
func (l customLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
}

// Error print error messages
func (l customLogger) Error(ctx context.Context, msg string, data ...interface{}) {
}

// Trace print sql message
func (l customLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	log.WithCtx(ctx).Info("database call", zap.String("sql", sql), zap.Int64("rows", rows), zap.String("elapsed", time.Since(begin).String()), zap.Error(err))
}
