package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
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
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,   // slow SQL threshold
				LogLevel:                  logger.Silent, // log level
				IgnoreRecordNotFoundError: true,          // ignore ErrRecordNotFound
				Colorful:                  false,         // disable color
			},
		),
	})
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return nil, fmt.Errorf("otelgorm plugin: %w", err)
	}
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
