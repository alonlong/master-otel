package models

import (
	commonv1 "master-otel/internal/proto/common/v1"
)

type User struct {
	ID       uint32 `gorm:"column:id"`
	Email    string `gorm:"column:email"`
	Username string `gorm:"column:username"`
}

func (e *User) ToProto() *commonv1.User {
	return &commonv1.User{
		Id:       e.ID,
		Email:    e.Email,
		Username: e.Username,
	}
}
