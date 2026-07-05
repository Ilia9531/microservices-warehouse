package payment

import def "Jopa/payment/internal/service"

var _ def.PaymentService = (*Service)(nil)

type Service struct{}

func NewService() *Service {
	return &Service{}
}
