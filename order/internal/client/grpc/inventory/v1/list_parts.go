package inventory

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	"github.com/Ilia9531/microservices-warehouse/order/internal/client/converter"
	"github.com/Ilia9531/microservices-warehouse/order/internal/model"
	grpcAuth "github.com/Ilia9531/microservices-warehouse/platform/pkg/middleware/grpc"
	inventoryv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/inventory/v1"
)

func (c *grpcClient) ListPartsByUUIDs(ctx context.Context, uuids []string) ([]*model.Part, error) {
	sessionUUID, ok := grpcAuth.GetSessionUUIDFromContext(ctx)
	fmt.Printf("📤 [INVENTORY-CLIENT] До Forward: session=%s, ok=%v\n", sessionUUID, ok)

	// превращает session_uuid из context в HTTP/2 заголовок gRPC
	ctx = grpcAuth.ForwardSessionUUIDToGRPC(ctx)

	sessionUUID2, ok := grpcAuth.GetSessionUUIDFromContext(ctx)
	fmt.Printf("📤 [INVENTORY-CLIENT] После Forward: session=%s, ok=%v\n", sessionUUID2, ok)

	// 🔍 2. ОТЛАДКА: проверяем, что metadata действительно добавилось
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		fmt.Printf("📤 [INVENTORY-CLIENT] Outgoing metadata: session-uuid=%v\n", md.Get("session-uuid"))
	} else {
		fmt.Println("⚠️ [INVENTORY-CLIENT] No outgoing metadata found")
	}

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
