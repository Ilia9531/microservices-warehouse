package v1

import (
	"context"

	"github.com/Ilia9531/microservices-warehouse/order/internal/converter"
	grpcAuth "github.com/Ilia9531/microservices-warehouse/platform/pkg/middleware/grpc"
	orderv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/openapi/order/v1"
	commonV1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/common/v1"
)

func (a *Api) APIV1OrdersPost(ctx context.Context,
	req *orderv1.CreateOrderRequest,
) (orderv1.APIV1OrdersPostRes, error) {
	// Читаем доверенный user_uuid из контекста (положен middleware'ом)
	user, ok := ctx.Value(grpcAuth.GetUserContextKey()).(*commonV1.User)

	if !ok || user.Uuid == "" {
		// Middleware не отработал → возвращаем 401
		return &orderv1.Error{
			Code:    orderv1.NewOptString("UNAUTHORIZED"),
			Message: orderv1.NewOptString("Authentication required"),
		}, nil
	}
	if req.PartUuids == nil {
		return &orderv1.Error{
			Code:    orderv1.NewOptString("VALIDATION_ERROR"),
			Message: orderv1.NewOptString("part_uuids cannot be empty"),
		}, nil
	}

	RespOrder, err := a.orderService.CreateOrder(ctx, user.Uuid, req.PartUuids)
	if err != nil {
		return nil, err
	}
	respProto := converter.ToCreateOrderResponse(RespOrder)

	return &orderv1.CreateOrderResponse{OrderUUID: respProto.OrderUUID, TotalPrice: respProto.TotalPrice}, nil
}
