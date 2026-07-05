package v1

import (
	"github.com/Ilia9531/microservices-warehouse/order/internal/model"
	orderv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/openapi/order/v1"
	"context"
	"errors"
)

func (a *Api) APIV1OrdersOrderUUIDCancelPost(ctx context.Context, params orderv1.APIV1OrdersOrderUUIDCancelPostParams) (orderv1.APIV1OrdersOrderUUIDCancelPostRes, error) {
	err := a.orderService.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		return a.handleCancelError(err)
	}
	return &orderv1.APIV1OrdersOrderUUIDCancelPostNoContent{}, nil
}

func (a *Api) handleCancelError(err error) (orderv1.APIV1OrdersOrderUUIDCancelPostRes, error) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		return &orderv1.APIV1OrdersOrderUUIDCancelPostNotFound{}, nil
	case errors.Is(err, model.ErrConflict):
		return &orderv1.Error{
			Code:    orderv1.NewOptString("CONFLICT"),
			Message: orderv1.NewOptString(err.Error()),
		}, nil
	case errors.Is(err, model.ErrValidation):
		return &orderv1.Error{
			Code:    orderv1.NewOptString("VALIDATION_ERROR"),
			Message: orderv1.NewOptString(err.Error()),
		}, nil
	default:
		return &orderv1.Error{
			Code:    orderv1.NewOptString("INTERNAL_SERVER_ERROR"),
			Message: orderv1.NewOptString("Internal server error"),
		}, nil
	}
}
