package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"Jopa/assembly/internal/model"
	"Jopa/platform/pkg/kafka"
	"Jopa/platform/pkg/logger"
	eventsv1 "Jopa/shared/pkg/proto/events/v1"
)

type service struct {
	producer kafka.Producer
}

func NewService(producer kafka.Producer) *service {
	return &service{producer: producer}
}

func (s *service) ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	pb := &eventsv1.ShipAssembledEvent{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(pb)
	if err != nil {
		logger.Error(ctx, "failed to marshal ShipAssembled", zap.Error(err))
		return err
	}

	if err := s.producer.Send(ctx, []byte(event.EventUUID), payload); err != nil {
		logger.Error(ctx, "failed to publish ShipAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "ShipAssembled published", zap.String("order_uuid", event.OrderUUID))
	return nil
}
