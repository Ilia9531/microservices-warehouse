package payment_test

import (
	"Jopa/payment/internal/model"

	"github.com/google/uuid"
)

// TestPayOrder_Success тест на успешную оплату.
func (s *PaymentServiceTestSuite) TestPayOrder_Success() {
	ctx := s.ctx

	// Вызываем сервис
	txUUID, err := s.service.PayOrder(ctx, "order-123", "user-456", "CARD")

	// Проверяем результат
	s.Require().NoError(err)
	s.Require().NotEmpty(txUUID)

	// Проверяем, что UUID валидный
	_, err = uuid.Parse(txUUID)
	s.Require().NoError(err, "transaction_uuid должен быть валидным UUID")
}

// TestPayOrder_Empty тест на пустые аргументы.
func (s *PaymentServiceTestSuite) TestPayOrder_Empty() {
	ctx := s.ctx
	// тест на пустой order_uuid
	txUUID, err := s.service.PayOrder(ctx, "", "user-456", "CARD")
	s.ErrorIs(err, model.ErrInvalidOrderUUID)
	s.Require().Empty(txUUID)

	// тест на пустой user_uuid
	txUUID, err = s.service.PayOrder(ctx, "order-121", "", "CARD")
	s.ErrorIs(err, model.ErrInvalidUserUUID)
	s.Require().Empty(txUUID)

	// тест на пустой paymentMethod
	txUUID, err = s.service.PayOrder(ctx, "order-121", "user-456", "")
	s.ErrorIs(err, model.ErrInvalidPaymentMethod)
	s.Require().Empty(txUUID)
}
