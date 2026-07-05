package order

import (
	"Jopa/order/internal/model"
	"Jopa/order/internal/repository/converter"
	"context"
	"fmt"
)

func (r *Repository) Update(ctx context.Context, orderM *model.Order) error {
	if orderM == nil || orderM.OrderUUID == "" {
		return fmt.Errorf("invalid order: order or order UUID is nil/empty")
	}
	RepoOrder := converter.ToRepo(orderM)

	//Здесь надо по хорошему for update использовать
	//внутри транзакции
	qyery := `UPDATE orders
	SET transaction_uuid = $1,
	    payment_method = $2,
	    status = $3,
	    updated_at = $4
	WHERE order_uuid = $5
	`

	_, err := r.db.Exec(ctx, qyery,
		RepoOrder.TransactionUUID,
		RepoOrder.PaymentMethod,
		RepoOrder.Status,
		RepoOrder.UpdatedAt,
		RepoOrder.OrderUUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update order in db %w", err)
	}
	return nil
}
