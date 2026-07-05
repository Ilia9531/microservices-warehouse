package order

import (
	"context"
	"fmt"

	"github.com/Ilia9531/microservices-warehouse/order/internal/model"
)

func (s *Service) GetOrder(ctx context.Context, uuid string) (*model.Order, error) {
	if uuid == "" {
		return nil, fmt.Errorf("uuid is required: %w", model.ErrValidation)
	}
	order, err := s.repo.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return order, nil
}
