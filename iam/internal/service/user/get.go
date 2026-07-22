package user

import (
	"context"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/model"
)

func (s *Service) GetUser(ctx context.Context, uuid string) (*model.User, error) {
	if uuid == "" {
		return nil, model.ErrInvalidLogin // или отдельная ошибка ErrEmptyUUID
	}
	return s.userRepo.GetByUUID(ctx, uuid)
}
