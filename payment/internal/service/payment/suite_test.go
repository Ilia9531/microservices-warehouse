package payment_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Ilia9531/microservices-warehouse/payment/internal/service/payment"
)

type PaymentServiceTestSuite struct {
	suite.Suite
	ctx     context.Context
	service *payment.Service
}

func (s *PaymentServiceTestSuite) SetupSuite() {
	s.service = payment.NewService()
}

// TestPaymentServiceSuite — точка входа для запуска всех тестов.
func TestPaymentServiceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(PaymentServiceTestSuite))
}
