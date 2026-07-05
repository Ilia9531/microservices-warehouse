package v1

import (
	"Jopa/inventory/internal/service"
	inventoryV1 "Jopa/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer
	InventoryService service.InventoryService
}

func NewApi(InventoryService service.InventoryService) *api {
	return &api{
		InventoryService: InventoryService,
	}
}
