package service

import (
	"context"

	"github.com/Ilia9531/microservices-warehouse/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (*model.Order, error)
	GetOrder(ctx context.Context, uuid string) (*model.Order, error)
	PayOrder(ctx context.Context, uuid, paymentMethod string) (*model.Order, error)
	CancelOrder(ctx context.Context, uuid string) error
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderProducerService interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error
}
