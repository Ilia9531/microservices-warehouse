package order_consumer

import (
	"context"
	"math/rand"
	"time"

	"go.uber.org/zap"

	"Jopa/assembly/internal/model"
	kafka "Jopa/platform/pkg/kafka/consumer"
	"Jopa/platform/pkg/logger"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Received OrderPaid, starting assembly",
		zap.String("order_uuid", event.OrderUUID),
		zap.String("topic", msg.Topic),
		zap.Int64("offset", msg.Offset),
	)

	// Имитация сборки: 1-10 секунд
	buildTime := int64(rand.Intn(10) + 1)
	select {
	case <-time.After(time.Duration(buildTime) * time.Second):
		// Сборка завершена
	case <-ctx.Done():
		logger.Warn(ctx, "Assembly interrupted by shutdown", zap.String("order_uuid", event.OrderUUID))
		return ctx.Err()
	}

	shipEvent := model.ShipAssembledEvent{
		EventUUID:    event.EventUUID,
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: buildTime,
	}

	if err := s.producerService.ProduceShipAssembled(ctx, shipEvent); err != nil {
		logger.Error(ctx, "Failed to publish ShipAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "ShipAssembled published",
		zap.String("order_uuid", event.OrderUUID),
		zap.Int64("build_time_sec", buildTime),
	)
	return nil
}
