package part

import (
	"Jopa/inventory/internal/model"
	"context"
	"fmt"
)

func (s *Service) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	if uuid == "" {
		return nil, fmt.Errorf("uuid is required: %w", model.ErrValidation)
	}
	part, err := s.inventoryRepository.Get(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("get part: %w", model.ErrPartNotFound)
	}
	return part, nil
}
