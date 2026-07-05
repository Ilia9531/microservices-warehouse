package part

import (
	"github.com/Ilia9531/microservices-warehouse/inventory/internal/repository"
	def "github.com/Ilia9531/microservices-warehouse/inventory/internal/service"
)

var _ def.InventoryService = (*Service)(nil)

type Service struct {
	inventoryRepository repository.InventoryRepository
}

func NewService(inventoryRepository repository.InventoryRepository) *Service {
	return &Service{
		inventoryRepository: inventoryRepository,
	}
}
