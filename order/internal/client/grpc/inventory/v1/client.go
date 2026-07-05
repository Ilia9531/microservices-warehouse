package inventory

import (
	inventoryv1 "Jopa/shared/pkg/proto/inventory/v1"

	cl "Jopa/order/internal/client/grpc"

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
