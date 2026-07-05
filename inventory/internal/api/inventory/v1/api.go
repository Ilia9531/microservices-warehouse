package v1

import (
	"github.com/Ilia9531/microservices-warehouse/inventory/internal/service"
	inventoryV1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/inventory/v1"
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
