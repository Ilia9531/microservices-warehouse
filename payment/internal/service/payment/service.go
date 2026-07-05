package payment

import def "github.com/Ilia9531/microservices-warehouse/payment/internal/service"

var _ def.PaymentService = (*Service)(nil)

type Service struct{}

func NewService() *Service {
	return &Service{}
}
