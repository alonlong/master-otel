package models

import (
	commonv1 "master-otel/internal/proto/common/v1"
)

type Email struct {
	Id      uint32 `gorm:"column:id;primary_key"`
	Email   string `gorm:"column:email"`
	Subject string `gorm:"column:subject"`
	Body    string `gorm:"column:body"`
}

func (e *Email) ToProto() *commonv1.Email {
	return &commonv1.Email{
		Id:      e.Id,
		Email:   e.Email,
		Subject: e.Subject,
		Body:    e.Body,
	}
}
