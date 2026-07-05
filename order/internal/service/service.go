package service

import (
	"Jopa/order/internal/model"
	"context"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (*model.Order, error)
	GetOrder(ctx context.Context, uuid string) (*model.Order, error)
	PayOrder(ctx context.Context, uuid string, paymentMethod string) (*model.Order, error)
	CancelOrder(ctx context.Context, uuid string) error
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderProducerService interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error
}
