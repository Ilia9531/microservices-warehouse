package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/model"
)

func (s *Service) Register(ctx context.Context, req *model.RegisterRequest) (string, error) {
	if err := req.Validate(); err != nil {
		return "", fmt.Errorf("validate RegisterRequest failed in service layer: %w", err)
	}

	// Проверяем, не существует ли пользователь с таким логином
	if _, err := s.userRepo.GetByLogin(ctx, req.Login); err == nil {
		return "", model.ErrUserAlreadyExists
	}

	// Хешируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Создаём доменную модель
	u := &model.User{
		UUID:                uuid.New().String(),
		Login:               req.Login,
		PasswordHash:        string(hash),
		Email:               req.Email,
		NotificationMethods: req.NotificationMethods,
	}

	// Сохраняем в репозиторий
	if err = s.userRepo.Create(ctx, u); err != nil {
		return "", err
	}

	return u.UUID, nil
}
