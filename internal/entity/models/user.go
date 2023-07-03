package models

import (
	storedv1 "master-otel/internal/proto/stored/v1"
)

type User struct {
	ID       uint32 `gorm:"column:id"`
	Email    string `gorm:"column:email"`
	Username string `gorm:"column:username"`
}

func (e *User) ToProto() *storedv1.User {
	return &storedv1.User{
		Email:    e.Email,
		Username: e.Username,
	}
}
