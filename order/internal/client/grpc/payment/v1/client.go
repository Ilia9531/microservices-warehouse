package payment

import (
	cl "github.com/Ilia9531/microservices-warehouse/order/internal/client/grpc"
	paymentv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
)

type grpcClient struct {
	client paymentv1.PaymentServiceClient
}

func NewPaymentClient(conn *grpc.ClientConn) cl.PaymentClient {
	return &grpcClient{
		client: paymentv1.NewPaymentServiceClient(conn),
	}
}
