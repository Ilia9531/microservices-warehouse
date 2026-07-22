package grpc

import (
	"context"

	"github.com/Ilia9531/microservices-warehouse/order/internal/model"
)

type InventoryClient interface {
	ListPartsByUUIDs(ctx context.Context, uuids []string) ([]*model.Part, error)
}

type PaymentClient interface {
	PayOrderCl(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error)
}

//type IamClient interface {
//	Whoami(ctx context.Context, sessionUUID string) (*model.Whoami, error)
//	Login(ctx context.Context, login, password string) (string, error)
//}
