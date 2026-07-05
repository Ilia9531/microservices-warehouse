package order_test

import (
	"github.com/Ilia9531/microservices-warehouse/order/internal/model"

	"github.com/stretchr/testify/mock"
)

// TestCancelOrder_Success тест на успешную отмену.
func (s *OrderServiceTestSuite) TestCancelOrder_Success() {
	ctx := s.ctx
	order := s.createTestOrder()

	s.repo.EXPECT().
		Get(ctx, "test-order-uuid").
		Return(order, nil).
		Once()

	s.repo.EXPECT().
		Update(ctx, mock.MatchedBy(func(o *model.Order) bool {
			return o.Status == model.StatusCancelled
		})).
		Return(nil).
		Once()
	err := s.service.CancelOrder(ctx, "test-order-uuid")

	s.Require().NoError(err)
}

// TestCancelOrder_AlreadyPaid тест на отмену оплаченного заказа.
func (s *OrderServiceTestSuite) TestCancelOrder_AlreadyPaid() {
	ctx := s.ctx
	order := s.createTestOrder()
	order.Status = model.StatusPaid

	s.repo.EXPECT().
		Get(ctx, "test-order-uuid").
		Return(order, nil).
		Once()

	err := s.service.CancelOrder(ctx, "test-order-uuid")

	s.Require().Error(err)
	s.ErrorIs(err, model.ErrConflict)
}

// TestCancelOrder_AlreadyCancelled тест на повторную отмену.
func (s *OrderServiceTestSuite) TestCancelOrder_AlreadyCancelled() {
	ctx := s.ctx
	order := s.createTestOrder()
	order.Status = model.StatusCancelled

	s.repo.EXPECT().
		Get(ctx, "test-order-uuid").
		Return(order, nil).
		Once()

	err := s.service.CancelOrder(ctx, "test-order-uuid")

	s.Require().Error(err)
	s.ErrorIs(err, model.ErrConflict)
}

// TestCancelOrder_NotFound тест на отсутствие заказа.
func (s *OrderServiceTestSuite) TestCancelOrder_NotFound() {
	ctx := s.ctx

	s.repo.EXPECT().
		Get(ctx, "non-existent").
		Return(nil, model.ErrNotFound).
		Once()

	err := s.service.CancelOrder(ctx, "non-existent")

	s.Require().Error(err)
	s.ErrorIs(err, model.ErrNotFound)
}

// TestCancelOrder_EmptyUUID тест на пустой UUID.
func (s *OrderServiceTestSuite) TestCancelOrder_EmptyUUID() {
	ctx := s.ctx

	err := s.service.CancelOrder(ctx, "")

	s.Require().Error(err)
	s.ErrorIs(err, model.ErrValidation)
}
