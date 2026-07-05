package part

import (
	"Jopa/inventory/internal/model"
	"context"
	"fmt"
)

func (s *Service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	parts, err := s.inventoryRepository.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get part: %w", model.ErrPartNotFound)
	}
	return parts, nil
}
