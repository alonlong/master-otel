package entity

import (
	"errors"

	"master-otel/internal/entity/models"

	"gorm.io/gorm"
)

const (
	UserTable = "users"
)

func (s *Store) CreateUser(entity *models.User) error {
	return s.db.Table(UserTable).Create(entity).Error
}

func (s *Store) GetUser(id uint32) (*models.User, error) {
	var entity models.User
	if err := s.db.Table(UserTable).Where("id = ?", id).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (s *Store) DeleteUser(id uint32) error {
	return s.db.Table(UserTable).Where("id = ?", id).Delete(&models.User{}).Error
}
