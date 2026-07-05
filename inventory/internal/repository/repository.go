package repository

import (
	"Jopa/inventory/internal/model"
	"context"
)

type InventoryRepository interface {
	Get(ctx context.Context, uuid string) (*model.Part, error)
	List(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error)
}
