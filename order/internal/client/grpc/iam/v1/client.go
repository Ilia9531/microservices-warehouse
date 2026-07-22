package iam

import (
	"context"

	"google.golang.org/grpc"

	authv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/auth/v1"
)

type grpcClient struct {
	authClient authv1.AuthServiceClient
}

// NewClient создаёт новый экземпляр gRPC-клиента к IAM сервису.

func NewClient(conn *grpc.ClientConn) *grpcClient {
	return &grpcClient{
		authClient: authv1.NewAuthServiceClient(conn),
	}
}

func (c *grpcClient) Whoami(ctx context.Context, sessionUUID string) (*authv1.WhoamiResponse, error) {
	req := &authv1.WhoamiRequest{
		SessionUuid: sessionUUID,
	}
	resp, err := c.authClient.Whoami(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//func (c *grpcClient) Login(_ context.Context, login, password string) (string, error) {
//	return fmt.Sprint("ops, method Logis wasn't creature:%s,%s\n", login, password), nil
//}
