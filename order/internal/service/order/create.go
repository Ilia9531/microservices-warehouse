package order

import (
	"github.com/Ilia9531/microservices-warehouse/order/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *Service) CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (*model.Order, error) {
	if userUUID == "" {
		return nil, fmt.Errorf("userUUID is empty: %w", model.ErrValidation)
	}
	if len(partUUIDs) == 0 {
		return nil, fmt.Errorf("partUUIDs is empty: %w", model.ErrValidation)
	}

	parts, err := s.inventoryClient.ListPartsByUUIDs(ctx, partUUIDs)

	if err != nil {
		return nil, fmt.Errorf("failed to list parts: %w", err)
	}
	if len(parts) != len(partUUIDs) {
		return nil, model.ErrNotFound
	}

	totalPrice := 0.0
	for _, v := range parts {
		totalPrice += v.Price
	}
	orderUUid := uuid.New().String()

	newOrder := &model.Order{
		OrderUUID:  orderUUid,
		UserUUID:   userUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = s.repo.Create(ctx, newOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", model.ErrConflict)
	}
	return newOrder, nil
}
