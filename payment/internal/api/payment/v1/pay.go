package v1

import (
	"github.com/Ilia9531/microservices-warehouse/payment/internal/model"
	paymentV1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/payment/v1"
	"context"

	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transUuid, err := a.service.PayOrder(
		ctx,
		req.OrderUuid,
		req.UserUuid,
		req.PaymentMethod.String(),
	)

	if err != nil {
		if errors.Is(err, model.ErrInvalidOrderUUID) {
			return nil, status.Errorf(codes.InvalidArgument, "Order Uuid is empty")
		}
		if errors.Is(err, model.ErrInvalidUserUUID) {
			return nil, status.Errorf(codes.InvalidArgument, "User Uuid is empty")
		}
		if errors.Is(err, model.ErrInvalidPaymentMethod) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid payment method")
		}
		return nil, err
	}
	return &paymentV1.PayOrderResponse{
		TransactionUuid: transUuid,
	}, nil
}
