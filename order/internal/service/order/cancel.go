package order

import (
	"Jopa/order/internal/model"
	"context"
	"fmt"
	"time"
)

func (s *Service) CancelOrder(ctx context.Context, uuid string) error {
	if uuid == "" {
		return fmt.Errorf("uuid is empty: %w", model.ErrValidation)
	}
	order, err := s.repo.Get(ctx, uuid)
	if err != nil {
		return err
	}
	if order.Status == model.StatusPaid {
		return fmt.Errorf("cannot cancel paid order: %w", model.ErrConflict)
	}
	if order.Status == model.StatusCancelled {
		return fmt.Errorf("cannot cancel cancelled order: %w", model.ErrConflict)
	}
	order.Status = model.StatusCancelled
	order.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, order); err != nil {
		return fmt.Errorf("cannot update order: %w", err)
	}
	return nil
}
