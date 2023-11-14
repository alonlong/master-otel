package entity

import (
	"context"
	"errors"
	"master-otel/internal/entity/models"

	"gorm.io/gorm"
)

const (
	EmailTable = "emails"
)

func (s *Store) CreateEmail(ctx context.Context, entity *models.Email) error {
	return s.db.WithContext(ctx).Table(EmailTable).Create(entity).Error
}

func (s *Store) GetEmail(ctx context.Context, id uint32) (*models.Email, error) {
	var entity models.Email
	if err := s.db.WithContext(ctx).Table(EmailTable).Where("id = ?", id).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (s *Store) DeleteEmail(ctx context.Context, id uint32) error {
	return s.db.WithContext(ctx).Table(EmailTable).Where("id = ?", id).Delete(&models.Email{}).Error
}
