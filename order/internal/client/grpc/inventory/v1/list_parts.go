package inventory

import (
	"context"

	"github.com/Ilia9531/microservices-warehouse/order/internal/client/converter"
	"github.com/Ilia9531/microservices-warehouse/order/internal/model"
	inventoryv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/inventory/v1"
)

func (c *grpcClient) ListPartsByUUIDs(ctx context.Context, uuids []string) ([]*model.Part, error) {
	resp, err := c.client.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{
			Uuids: uuids,
		},
	})
	if err != nil {
		return nil, err
	}
	parts := converter.ToDomainParts(resp.Parts)
	return parts, nil
}
