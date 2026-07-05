package service

import (
	"github.com/Ilia9531/microservices-warehouse/inventory/internal/model"
	"context"
)

type InventoryService interface {
	GetPart(ctx context.Context, uuid string) (*model.Part, error)
	ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error)
}
