package models

import (
	"time"

	commonv1 "master-otel/internal/proto/common/v1"
)

type Attachment struct {
	Id         uint32    `gorm:"column:id;primary_key"`
	Email      uint32    `gorm:"column:email"`
	Attachment string    `gorm:"column:attachment"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (a *Attachment) ToProto() *commonv1.Attachment {
	return &commonv1.Attachment{
		Id:         a.Id,
		Email:      a.Email,
		Attachment: a.Attachment,
	}
}
