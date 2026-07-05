package order_test

import (
	"testing"

	"context"

	"github.com/stretchr/testify/suite"

	clmocks "Jopa/order/internal/client/grpc/mocks"
	"Jopa/order/internal/model"
	repomocks "Jopa/order/internal/repository/mocks"
	servmocks "Jopa/order/internal/service/mocks"
	"Jopa/order/internal/service/order"
)

// OrderServiceTestSuite - тестовый набор для OrderService.
type OrderServiceTestSuite struct {
	suite.Suite
	ctx             context.Context
	repo            *repomocks.OrderRepository
	inventoryClient *clmocks.InventoryClient
	paymentClient   *clmocks.PaymentClient
	service         *order.Service
	producerService *servmocks.OrderProducerService
}

// SetupTest выполняется перед каждым тестом.
func (s *OrderServiceTestSuite) SetupTest() {
	// Создаём моки
	s.repo = repomocks.NewOrderRepository(s.T())
	s.inventoryClient = clmocks.NewInventoryClient(s.T())
	s.paymentClient = clmocks.NewPaymentClient(s.T())
	s.producerService = servmocks.NewOrderProducerService(s.T())

	// Создаём сервис с моками
	s.service = order.NewService(s.repo, s.inventoryClient, s.paymentClient, s.producerService)
}

// TestOrderServiceSuite запускает тестовый набор.
func TestOrderServiceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(OrderServiceTestSuite))
}

// Helper: создание тестового заказа.
func (s *OrderServiceTestSuite) createTestOrder() *model.Order {
	return &model.Order{
		OrderUUID:  "test-order-uuid",
		UserUUID:   "test-user-uuid",
		PartUUIDs:  []string{"part-1", "part-2"},
		TotalPrice: 199.99,
		Status:     model.StatusPendingPayment,
	}
}
