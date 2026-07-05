package payment_test

import (
	"Jopa/payment/internal/service/payment"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
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
