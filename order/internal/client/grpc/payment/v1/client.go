package payment

import (
	cl "Jopa/order/internal/client/grpc"
	paymentv1 "Jopa/shared/pkg/proto/payment/v1"
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
