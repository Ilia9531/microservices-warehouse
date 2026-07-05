package order_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Jopa/order/internal/model"
)

// TestPayOrder_Success тест на успешную оплату.
func (s *OrderServiceTestSuite) TestPayOrder_Success() {
	ctx := s.ctx
	order := s.createTestOrder()
	paymentMethod := "CARD"

	// Получение заказа
	s.repo.EXPECT().
		Get(ctx, "test-order-uuid").
		Return(order, nil).
		Once()

	// Вызов PaymentClient
	s.paymentClient.EXPECT().
		PayOrderCl(ctx, order.OrderUUID, order.UserUUID, paymentMethod).
		Return("transaction-uuid", nil).
		Once()

	// Обновление заказа
	s.repo.EXPECT().
		Update(mock.Anything, mock.MatchedBy(func(o *model.Order) bool {
			return o.Status == model.StatusPaid
		})).
		Return(nil).
		Once()

	result, err := s.service.PayOrder(ctx, "test-order-uuid", paymentMethod)

	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Equal(model.StatusPaid, result.Status)
	s.Equal("transaction-uuid", *result.TransactionUUID)
}

// TestPayOrder_AlreadyPaid тест на повторную оплату.
func (s *OrderServiceTestSuite) TestPayOrder_AlreadyPaid() {
	ctx := s.ctx
	order := s.createTestOrder()
	order.Status = model.StatusPaid
	paymentMethod := "CARD"

	s.repo.EXPECT().
		Get(ctx, "test-order-uuid").
		Return(order, nil).
		Once()

	result, err := s.service.PayOrder(ctx, "test-order-uuid", paymentMethod)

	s.Require().Error(err)
	s.Require().Nil(result)
	s.ErrorIs(err, model.ErrConflict)
}

// TestPayOrder_Cancelled тест на оплату отменённого заказа.
func (s *OrderServiceTestSuite) TestPayOrder_Cancelled() {
	ctx := s.ctx
	order := s.createTestOrder()
	order.Status = model.StatusCancelled
	paymentMethod := "CARD"

	s.repo.EXPECT().
		Get(ctx, "test-order-uuid").
		Return(order, nil).
		Once()

	result, err := s.service.PayOrder(ctx, "test-order-uuid", paymentMethod)

	s.Require().Error(err)
	s.Require().Nil(result)
	s.ErrorIs(err, model.ErrConflict)
}

// TestPayOrder_NotFound тест на отсутствие заказа.
func (s *OrderServiceTestSuite) TestPayOrder_NotFound() {
	ctx := s.ctx
	paymentMethod := "CARD"

	s.repo.EXPECT().
		Get(ctx, "non-existent").
		Return(nil, model.ErrNotFound).
		Once()

	result, err := s.service.PayOrder(ctx, "non-existent", paymentMethod)

	s.Require().Error(err)
	s.Require().Nil(result)
	s.ErrorIs(err, model.ErrNotFound)
}

// TestPayOrder_EmptyUUID тест на пустой UUID.
func (s *OrderServiceTestSuite) TestPayOrder_EmptyUUID() {
	ctx := s.ctx
	paymentMethod := "CARD"

	result, err := s.service.PayOrder(ctx, "", paymentMethod)

	s.Require().Error(err)
	s.Require().Nil(result)
	s.ErrorIs(err, model.ErrValidation)
}

// TestPayOrder_EmptyPaymentMethod тест на пустой paymentMethod.
func (s *OrderServiceTestSuite) TestPayOrder_EmptyPaymentMethod() {
	ctx := s.ctx
	paymentMethod := ""

	result, err := s.service.PayOrder(ctx, "test-uuid", paymentMethod)

	s.Require().Error(err)
	s.Require().Nil(result)
	s.ErrorIs(err, model.ErrValidation)
}

// TestPayOrder_PaymentError тест на ошибку PaymentClient.
func (s *OrderServiceTestSuite) TestPayOrder_PaymentError() {
	ctx := s.ctx
	order := s.createTestOrder()
	paymentMethod := "CARD"

	s.repo.EXPECT().
		Get(ctx, "test-order-uuid").
		Return(order, nil).
		Once()

	s.paymentClient.EXPECT().
		PayOrderCl(ctx, order.OrderUUID, order.UserUUID, paymentMethod).
		Return("", assert.AnError).
		Once()

	result, err := s.service.PayOrder(ctx, "test-order-uuid", paymentMethod)

	s.Require().Error(err)
	s.Require().Nil(result)
}
