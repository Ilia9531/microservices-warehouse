package v1

import (
	"github.com/Ilia9531/microservices-warehouse/payment/internal/service"
	paymentV1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer
	service service.PaymentService
}

func NewApi(service service.PaymentService) *api {
	return &api{
		service: service,
	}
}
