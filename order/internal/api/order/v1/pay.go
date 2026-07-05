package v1

import (
	"github.com/Ilia9531/microservices-warehouse/order/internal/converter"
	orderv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/openapi/order/v1"
	"context"
	"log"
)

func (a *Api) APIV1OrdersOrderUUIDPayPost(ctx context.Context, req *orderv1.PayOrderRequest,
	params orderv1.APIV1OrdersOrderUUIDPayPostParams) (orderv1.APIV1OrdersOrderUUIDPayPostRes, error) {
	if params.OrderUUID == "" || req.PaymentMethod == "" {
		log.Printf("Ошибка при проверке OrderUUID и PaymentMethod")
		return &orderv1.Error{
			Code:    orderv1.NewOptString("VALIDATION_ERROR"),
			Message: orderv1.NewOptString("OrderUUID and PaymentMethod cannot be empty"),
		}, nil
	}
	paymentMethod := string(req.PaymentMethod)
	respOrder, err := a.orderService.PayOrder(ctx, params.OrderUUID, paymentMethod)
	if err != nil {
		return nil, err
	}
	transactionProto := converter.ToPayOrderResponse(respOrder.TransactionUUID)
	return transactionProto, nil
}
