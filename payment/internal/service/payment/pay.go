package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/Ilia9531/microservices-warehouse/payment/internal/model"
)

func (s *Service) PayOrder(_ context.Context, orderUUID, userUUID, paymentMethod string) (string, error) {
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

	slog.Info("Оплата прошла успешно", "transaction_uuid", transUuid)
	return transUuid, nil
}
