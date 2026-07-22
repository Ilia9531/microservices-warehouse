package model

import (
	"errors"

	common "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/common/v1"
)

type User struct {
	UUID                string
	Login               string
	PasswordHash        string
	Email               string
	NotificationMethods []*common.NotificationMethod
}

type RegisterRequest struct {
	Login               string
	Password            string
	Email               string
	NotificationMethods []*common.NotificationMethod
}

// Validate проверяет корректность данных для регистрации.
func (r *RegisterRequest) Validate() error {
	if r.Login == "" {
		return ErrInvalidLogin
	}
	if r.Password == "" {
		return ErrInvalidPassword
	}
	if r.Email == "" {
		return ErrInvalidEmail
	}
	return nil
}

// Validate проверяет базовые инварианты пользователя.
// Не включает проверку пароля — это задача сервиса.
func (u *User) Validate() error {
	if u.UUID == "" {
		return errors.New("user UUID is required")
	}
	if u.Login == "" {
		return ErrInvalidLogin
	}
	if u.Email == "" {
		return ErrInvalidEmail
	}
	return nil
}
