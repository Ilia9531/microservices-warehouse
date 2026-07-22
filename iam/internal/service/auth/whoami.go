package auth

import (
	"context"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/model"
)

func (s *Service) Whoami(ctx context.Context, sessionUUID string) (*model.Whoami, error) {
	if sessionUUID == "" {
		return nil, model.ErrSessionInvalid
	}

	// Получаем сессию из Redis
	session, err := s.sessionRepo.GetByUUID(ctx, sessionUUID)
	if err != nil {
		return nil, err
	}

	// Получаем пользователя по UUID из сессии
	user, err := s.userRepo.GetByUUID(ctx, session.UserUUID)
	if err != nil {
		return nil, err
	}
	whoami := &model.Whoami{
		UUID:  user.UUID,
		Login: user.Login,
		Email: user.Email,
	}
	return whoami, nil
}
