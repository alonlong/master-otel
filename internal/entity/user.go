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

func (s *Store) CreateUser(ctx context.Context, e *models.User) error {
	return s.db.WithContext(ctx).Table(User).Create(e).Error
}

func (s *Store) GetUser(ctx context.Context, id int64) (*models.User, error) {
	var e models.User
	if err := s.db.WithContext(ctx).Table(User).Where("id = ?", id).First(&e).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (s *Store) DeleteUser(ctx context.Context, id int64) error {
	return s.db.WithContext(ctx).Table(User).Where("id = ?", id).Delete(&models.User{}).Error
}
