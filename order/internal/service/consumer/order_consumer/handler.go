package order_consumer

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"Jopa/order/internal/model"
	kafka "Jopa/platform/pkg/kafka/consumer"
	"Jopa/platform/pkg/logger"
	"errors"
)

// ShipAssembledHandler обрабатывает входящее событие ShipAssembled из Kafka.
// Вызывается платформенным consumer для каждого сообщения в топике order.assembled.
func (s *service) ShipAssembledHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.shipAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembled event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing ShipAssembled event",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	// Здесь вызов бизнес-логики для обновления статуса заказа на ASSEMBLED.

	// Прямой вызов репозитория для обновления статуса
	ord, err := s.orderRepo.Get(ctx, event.OrderUUID)
	if err != nil {
		return fmt.Errorf("get order for assembly: %w", err)
	}
	if ord == nil {
		return errors.New("order not found for assembly")
	}
	if ord.Status != model.StatusPaid {
		logger.Warn(ctx, "Invalid status transition for assembly",
			zap.String("order_uuid", event.OrderUUID),
			zap.String("current_status", string(ord.Status)),
		)
		return fmt.Errorf("order must be in PAID status, got %s", ord.Status)
	}
	ord.Status = model.StatusAssembled
	if err = s.orderRepo.Update(ctx, ord); err != nil {
		return fmt.Errorf("update order status to ASSEMBLED: %w", err)
	}
	logger.Info(ctx, "Order status updated to ASSEMBLED",
		zap.String("order_uuid", event.OrderUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	return nil
}
