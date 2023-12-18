package stored

import (
	"context"
	"fmt"

	"master-otel/internal/entity"
	storedv1 "master-otel/internal/proto/stored/v1"

	"master-otel/pkg/db"
)

type Config struct {
	DB *db.Config
}

type Service struct {
	storedv1.UnimplementedStoredServiceServer

	store *entity.Store
}

func NewService(cfg *Config) (*Service, error) {
	store, err := entity.NewStore(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("new store: %w", err)
	}
	return &Service{
		store: store,
	}, nil
}

func (s *Service) Run(ctx context.Context) {}

func (s *Service) Shutdown() {
	if s.store != nil {
		s.store.Close()
	}
}
