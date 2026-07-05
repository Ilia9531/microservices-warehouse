package order

import (
	cl "github.com/Ilia9531/microservices-warehouse/order/internal/client/grpc"
	"github.com/Ilia9531/microservices-warehouse/order/internal/repository"
	"github.com/Ilia9531/microservices-warehouse/order/internal/service"
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
