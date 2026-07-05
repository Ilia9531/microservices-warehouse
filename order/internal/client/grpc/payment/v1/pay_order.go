package payment

import (
	"context"
	"fmt"

	"github.com/Ilia9531/microservices-warehouse/order/internal/converter"
	"github.com/Ilia9531/microservices-warehouse/order/internal/model"
	paymentv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/payment/v1"
)

func (c *grpcClient) PayOrderCl(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error) {
	protoPaymentMethod, err := converter.StrToMethod(paymentMethod)
	if err != nil {
		return "", fmt.Errorf("paymentMethod is invalid: %w", model.ErrValidation)
	}

	resp, err := c.client.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: protoPaymentMethod,
	})
	if err != nil {
		return "", err
	}
	return resp.TransactionUuid, nil
}
