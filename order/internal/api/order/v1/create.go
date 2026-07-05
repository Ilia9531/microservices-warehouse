package v1

import (
	"context"

	"github.com/Ilia9531/microservices-warehouse/order/internal/converter"
	orderv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/openapi/order/v1"
)

func (a *Api) APIV1OrdersPost(ctx context.Context,
	req *orderv1.CreateOrderRequest,
) (orderv1.APIV1OrdersPostRes, error) {
	if req.PartUuids == nil || req.UserUUID == "" {
		return &orderv1.Error{
			Code:    orderv1.NewOptString("VALIDATION_ERROR"),
			Message: orderv1.NewOptString("part_uuids cannot be empty"),
		}, nil
	}
	RespOrder, err := a.orderService.CreateOrder(ctx, req.UserUUID, req.PartUuids)
	if err != nil {
		return nil, err
	}
	respProto := converter.ToCreateOrderResponse(RespOrder)

	return &orderv1.CreateOrderResponse{OrderUUID: respProto.OrderUUID, TotalPrice: respProto.TotalPrice}, nil
}
