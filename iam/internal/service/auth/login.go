package auth

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/model"
)

func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	if login == "" {
		return "", model.ErrInvalidLogin
	}
	if password == "" {
		return "", model.ErrInvalidPassword
	}

	// Находим пользователя по логину
	u, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		return "", err
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", model.ErrInvalidPassword
	}

	// Создаём сессию
	session := &model.Session{
		UUID:     uuid.New().String(),
		UserUUID: u.UUID,
	}

	// Сохраняем в Redis с TTL
	if err = s.sessionRepo.Create(ctx, session, s.sessionTTL); err != nil {
		return "", err
	}

	return session.UUID, nil
}
