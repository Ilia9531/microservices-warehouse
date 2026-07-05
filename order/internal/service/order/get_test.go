package order_test

import (
	"context"

	"Jopa/order/internal/model"
)

// TestGetOrder_Success тест на успешное получение заказа.
func (s *OrderServiceTestSuite) TestGetOrder_Success() {
	ctx := s.ctx
	expectedOrder := s.createTestOrder()

	s.repo.EXPECT().
		Get(ctx, "test-order-uuid").
		Return(expectedOrder, nil).
		Once()

	order, err := s.service.GetOrder(ctx, "test-order-uuid")

	s.Require().NoError(err)
	s.Require().NotNil(order)
	s.Equal("test-order-uuid", order.OrderUUID)
}

// TestGetOrder_NotFound тест на отсутствие заказа.
func (s *OrderServiceTestSuite) TestGetOrder_NotFound() {
	ctx := s.ctx

	s.repo.EXPECT().
		Get(ctx, "non-existent").
		Return(nil, model.ErrNotFound).
		Once()

	order, err := s.service.GetOrder(ctx, "non-existent")

	s.Require().Error(err)
	s.Require().Nil(order)
	s.ErrorIs(err, model.ErrNotFound)
}

// TestGetOrder_EmptyUUID тест на пустой UUID.
func (s *OrderServiceTestSuite) TestGetOrder_EmptyUUID() {
	ctx := context.Background()

	order, err := s.service.GetOrder(ctx, "")

	s.Require().Error(err)
	s.Require().Nil(order)
	s.ErrorIs(err, model.ErrValidation)
}
