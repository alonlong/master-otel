package entity

import (
	"fmt"

	"master-otel/pkg/db"
	"master-otel/pkg/log"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(cfg *db.Config) (*Store, error) {
	db, err := db.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("new db: %w", err)
	}
	log.Info("init postgres", zap.String("addr", cfg.Host+":"+cfg.Port))
	return &Store{db: db}, nil
}

func (s *Store) Close() {
	if s.db != nil {
		if raw, err := s.db.DB(); err == nil {
			if err := raw.Close(); err != nil {
				log.Error("db close", zap.Error(err))
			}
		}
	}
}
