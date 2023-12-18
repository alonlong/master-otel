package entity

import (
	"context"
	"errors"

	"master-otel/internal/entity/models"

	"gorm.io/gorm"
)

const (
	Attachment = "attachments"
)

func (s *Store) CreateAttachment(ctx context.Context, e *models.Attachment) error {
	return s.db.WithContext(ctx).Table(Attachment).Create(e).Error
}

func (s *Store) GetAttachment(ctx context.Context, id uint32) (*models.Attachment, error) {
	e := new(models.Attachment)
	if err := s.db.WithContext(ctx).Table(Attachment).Where("id = ?", id).First(e).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return e, nil
}

func (s *Store) DeleteAttachment(ctx context.Context, id uint32) error {
	return s.db.WithContext(ctx).Table(Attachment).Where("id = ?", id).Delete(&models.Attachment{}).Error
}
