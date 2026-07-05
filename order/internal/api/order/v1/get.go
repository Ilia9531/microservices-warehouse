package v1

import (
	"context"

	"github.com/Ilia9531/microservices-warehouse/order/internal/converter"
	orderv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/openapi/order/v1"
)

func (a *Api) APIV1OrdersOrderUUIDGet(
	ctx context.Context, params orderv1.APIV1OrdersOrderUUIDGetParams,
) (orderv1.APIV1OrdersOrderUUIDGetRes, error) {
	if params.OrderUUID == "" {
		return &orderv1.Error{
			Code:    orderv1.NewOptString("VALIDATION_ERROR"),
			Message: orderv1.NewOptString("part_uuids cannot be empty"),
		}, nil
	}
	respOrder, err := a.orderService.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		return &orderv1.APIV1OrdersOrderUUIDGetNotFound{}, nil
	}
	respProto := converter.ToGetOrderResponse(respOrder)
	return respProto, nil
}
