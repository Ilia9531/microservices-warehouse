package converter

import (
	"github.com/Ilia9531/microservices-warehouse/order/internal/model"
	repoModel "github.com/Ilia9531/microservices-warehouse/order/internal/repository/model"
)

func ToRepo(order *model.Order) *repoModel.Order {
	if order == nil {
		return nil
	}
	return &repoModel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          string(order.Status),
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

func ToDomain(repoOrder *repoModel.Order) *model.Order {
	if repoOrder == nil {
		return nil
	}
	return &model.Order{
		OrderUUID:       repoOrder.OrderUUID,
		UserUUID:        repoOrder.UserUUID,
		PartUUIDs:       repoOrder.PartUUIDs,
		TotalPrice:      repoOrder.TotalPrice,
		TransactionUUID: repoOrder.TransactionUUID,
		PaymentMethod:   repoOrder.PaymentMethod,
		Status:          model.OrderStatus(repoOrder.Status),
		CreatedAt:       repoOrder.CreatedAt,
		UpdatedAt:       repoOrder.UpdatedAt,
	}
}
