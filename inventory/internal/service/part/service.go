package part

import (
	"Jopa/inventory/internal/repository"
	def "Jopa/inventory/internal/service"
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
