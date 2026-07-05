package inventory

import (
	inventoryv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/inventory/v1"

	cl "github.com/Ilia9531/microservices-warehouse/order/internal/client/grpc"

	"google.golang.org/grpc"
)

// grpcClient implements InventoryClient.
type grpcClient struct {
	client inventoryv1.InventoryServiceClient
}

// NewClient creates a new InventoryClient.
func NewClient(conn *grpc.ClientConn) cl.InventoryClient {
	return &grpcClient{
		client: inventoryv1.NewInventoryServiceClient(conn),
	}
}
