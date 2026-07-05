package v1

import (
	"Jopa/order/internal/service"
)

type Api struct {
	orderService service.OrderService
}

func NewApi(OrderService service.OrderService) *Api {
	return &Api{OrderService}
}
