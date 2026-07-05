package inventory

import (
	"Jopa/order/internal/client/converter"
	"Jopa/order/internal/model"
	inventoryv1 "Jopa/shared/pkg/proto/inventory/v1"
	"context"
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
