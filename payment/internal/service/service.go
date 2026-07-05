package service

import (
	"context"
)

type PaymentService interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod string) (string, error)
}
