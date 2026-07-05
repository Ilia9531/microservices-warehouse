package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/Ilia9531/microservices-warehouse/order/internal/model"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/kafka"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/logger"
	eventsv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/events/v1"
)

type service struct {
	orderPaidProducer kafka.Producer
}

func NewService(orderPaidProducer kafka.Producer) *service {
	return &service{
		orderPaidProducer: orderPaidProducer,
	}
}

func (p *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	msg := &eventsv1.OrderPaidEvent{
		EventUuid:       event.EventUUID,
		OrderUuid:       event.OrderUUID,
		UserUuid:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionUUID,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderPaidEvent", zap.Error(err))
		return err
	}

	// event_uuid используется как ключ партиционирования — гарантирует, что дубликаты
	// одного события попадут в одну партицию и будут обработаны идемпотентно.
	err = p.orderPaidProducer.Send(ctx, []byte(event.EventUUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish OrderPaid", zap.Error(err))
		return err
	}
	//для ниже логер для отладки
	logger.Info(ctx, "OrderPaid event published",
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
	)

	return nil
}
