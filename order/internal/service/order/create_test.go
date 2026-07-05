package order_test

import (
	"Jopa/order/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestCreateOrder_Success тест на успешное создание заказа.
func (s *OrderServiceTestSuite) TestCreateOrder_Success() {
	ctx := s.ctx
	userUUID := "user-123"
	partUUIDs := []string{"part-1", "part-2"}

	// Мокируем ответ InventoryClient
	parts := []*model.Part{
		{UUID: "part-1", Price: 99.99},
		{UUID: "part-2", Price: 100.00},
	}
	s.inventoryClient.EXPECT().
		ListPartsByUUIDs(ctx, partUUIDs).
		Return(parts, nil).
		Once()

	// Мокируем успешное создание в репозитории
	s.repo.EXPECT().
		Create(ctx, mock.MatchedBy(func(o *model.Order) bool {
			return o.UserUUID == userUUID && len(o.PartUUIDs) == 2
		})).
		Return(nil).
		Once()

	// Вызываем сервис
	order, err := s.service.CreateOrder(ctx, userUUID, partUUIDs)

	// Проверяем результат
	s.Require().NoError(err)
	s.Require().NotNil(order)
	s.Equal(userUUID, order.UserUUID)
	s.Equal(199.99, order.TotalPrice)
	s.Equal(model.StatusPendingPayment, order.Status)
}

// TestCreateOrder_EmptyPartUUIDs тест на пустой список деталей.
func (s *OrderServiceTestSuite) TestCreateOrder_EmptyPartUUIDs() {
	ctx := s.ctx

	order, err := s.service.CreateOrder(ctx, "user-123", []string{})

	s.Require().Error(err)
	s.Require().Nil(order)
	s.ErrorIs(err, model.ErrValidation)
}

// TestCreateOrder_EmptyUserUUID тест на пустой userUUID.
func (s *OrderServiceTestSuite) TestCreateOrder_EmptyUserUUID() {
	ctx := s.ctx

	order, err := s.service.CreateOrder(ctx, "", []string{"part-1"})

	s.Require().Error(err)
	s.Require().Nil(order)
	s.ErrorIs(err, model.ErrValidation)
}

// TestCreateOrder_PartsNotFound тест на отсутствие деталей в Inventory.
func (s *OrderServiceTestSuite) TestCreateOrder_PartsNotFound() {
	ctx := s.ctx
	partUUIDs := []string{"part-1", "part-2"}

	// Возвращаем только 1 часть вместо 2
	parts := []*model.Part{
		{UUID: "part-1", Price: 99.99},
	}
	s.inventoryClient.EXPECT().
		ListPartsByUUIDs(ctx, partUUIDs).
		Return(parts, nil).
		Once()

	order, err := s.service.CreateOrder(ctx, "user-123", partUUIDs)

	s.Require().Error(err)
	s.Require().Nil(order)
	s.ErrorIs(err, model.ErrNotFound)
}

// TestCreateOrder_InventoryError тест на ошибку InventoryClient.
func (s *OrderServiceTestSuite) TestCreateOrder_InventoryError() {
	ctx := s.ctx

	s.inventoryClient.EXPECT().
		ListPartsByUUIDs(ctx, []string{"part-1"}).
		Return(nil, assert.AnError).
		Once()

	order, err := s.service.CreateOrder(ctx, "user-123", []string{"part-1"})

	s.Require().Error(err)
	s.Require().Nil(order)
}

// TestCreateOrder_RepositoryError тест на ошибку репозитория.
func (s *OrderServiceTestSuite) TestCreateOrder_RepositoryError() {
	ctx := s.ctx
	partUUIDs := []string{"part-1"}

	parts := []*model.Part{
		{UUID: "part-1", Price: 99.99},
	}
	s.inventoryClient.EXPECT().
		ListPartsByUUIDs(ctx, partUUIDs).
		Return(parts, nil).
		Once()

	s.repo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(assert.AnError).
		Once()

	order, err := s.service.CreateOrder(ctx, "user-123", partUUIDs)

	s.Require().Error(err)
	s.Require().Nil(order)
}
