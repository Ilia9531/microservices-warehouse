package grpc

import (
	"Jopa/order/internal/model"
	"context"
)

type InventoryClient interface {
	ListPartsByUUIDs(ctx context.Context, uuids []string) ([]*model.Part, error)
}

type PaymentClient interface {
	PayOrderCl(ctx context.Context, orderUUID, userUUID string, paymentMethod string) (string, error)
}
