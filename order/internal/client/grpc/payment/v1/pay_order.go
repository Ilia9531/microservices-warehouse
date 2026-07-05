package payment

import (
	"Jopa/order/internal/converter"
	"Jopa/order/internal/model"
	paymentv1 "Jopa/shared/pkg/proto/payment/v1"
	"context"
	"fmt"
)

func (c *grpcClient) PayOrderCl(ctx context.Context, orderUUID, userUUID string, paymentMethod string) (string, error) {
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
