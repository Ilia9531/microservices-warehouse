package order

import (
	"context"
	"fmt"

	Guuid "github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Ilia9531/microservices-warehouse/order/internal/model"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/logger"
)

func (s *Service) PayOrder(ctx context.Context, uuid, paymentMethod string) (*model.Order, error) {
	if uuid == "" {
		return nil, fmt.Errorf("uuid is required: %w", model.ErrValidation)
	}
	if paymentMethod == "" {
		return nil, fmt.Errorf("paymentMethod is required: %w", model.ErrValidation)
	}

	order, err := s.repo.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if order.Status == model.StatusPaid {
		return nil, fmt.Errorf("order is already paid: %w", model.ErrConflict)
	}
	if order.Status == model.StatusCancelled {
		return nil, fmt.Errorf("order is already cancelled: %w", model.ErrConflict)
	}

	txUUID, err := s.paymentClient.PayOrderCl(ctx, order.OrderUUID, order.UserUUID, paymentMethod)
	if err != nil {
		return nil, fmt.Errorf("payment failed: %w", err)
	}
	order.Status = model.StatusPaid
	order.TransactionUUID = &txUUID
	order.PaymentMethod = &paymentMethod
	err = s.repo.Update(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	event := model.OrderPaidEvent{
		EventUUID:       Guuid.New().String(),
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PaymentMethod:   paymentMethod,
		TransactionUUID: txUUID,
	}
	err = s.producerService.ProduceOrderPaid(ctx, event)
	if err != nil {
		logger.Warn(ctx, "failed to publish OrderPaid, will retry asynchronously", zap.Error(err))
	}
	return order, nil
}
