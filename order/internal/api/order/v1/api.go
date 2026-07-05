package v1

import (
	"github.com/Ilia9531/microservices-warehouse/order/internal/service"
)

type Api struct {
	orderService service.OrderService
}

func NewApi(OrderService service.OrderService) *Api {
	return &Api{OrderService}
}
