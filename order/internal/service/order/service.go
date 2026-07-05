package order

import (
	cl "Jopa/order/internal/client/grpc"
	"Jopa/order/internal/repository"
	"Jopa/order/internal/service"
)

var _ service.OrderService = (*Service)(nil)

type Service struct {
	repo            repository.OrderRepository
	inventoryClient cl.InventoryClient
	paymentClient   cl.PaymentClient
	producerService service.OrderProducerService
}

func NewService(
	repo repository.OrderRepository,
	invClient cl.InventoryClient,
	payClient cl.PaymentClient,
	producerService service.OrderProducerService,
) *Service {
	return &Service{
		repo:            repo,
		inventoryClient: invClient,
		paymentClient:   payClient,
		producerService: producerService,
	}
}
