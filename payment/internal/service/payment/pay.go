package payment

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Ilia9531/microservices-warehouse/payment/internal/model"

	"github.com/google/uuid"
)

func (s *Service) PayOrder(_ context.Context, orderUUID, userUUID string, paymentMethod string) (string, error) {
	if orderUUID == "" {
		return "", model.ErrInvalidOrderUUID
	}
	if userUUID == "" {
		return "", model.ErrInvalidUserUUID
	}
	if paymentMethod == "" {
		return "", model.ErrInvalidPaymentMethod
	}
	transUuid := uuid.New().String()
	slog.Info("Processing payment",
		"user_uuid", userUUID,
		"order_uuid", orderUUID,
		"payment_method", paymentMethod,
		"transaction_uuid", transUuid,
	)
	fmt.Printf("Оплата прошла успешно, transaction_uuid: %s", transUuid)
	return transUuid, nil
}
