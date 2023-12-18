package entity

import (
	"context"
	"errors"

	"master-otel/internal/entity/models"

	"gorm.io/gorm"
)

const (
	User = "users"
)

func (s *Store) CreateUser(ctx context.Context, entity *models.User) error {
	return s.db.WithContext(ctx).Table(User).Create(entity).Error
}

func (s *Store) GetUser(ctx context.Context, id uint32) (*models.User, error) {
	var entity models.User
	if err := s.db.WithContext(ctx).Table(User).Where("id = ?", id).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (s *Store) DeleteUser(ctx context.Context, id uint32) error {
	return s.db.WithContext(ctx).Table(User).Where("id = ?", id).Delete(&models.User{}).Error
}
