package order

import (
	"Jopa/order/internal/model"
	"Jopa/order/internal/repository/converter"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) Create(ctx context.Context, orderM *model.Order) error {
	if orderM == nil || orderM.OrderUUID == "" {
		return fmt.Errorf("invalid order: order or order UUID is nil/empty")
	}
	dbOrder := converter.ToRepo(orderM)

	query := `
		INSERT INTO orders(
		    order_uuid, user_uuid, part_uuids, total_price,
		    status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		dbOrder.OrderUUID,
		dbOrder.UserUUID,
		dbOrder.PartUUIDs,
		dbOrder.TotalPrice,
		dbOrder.Status,
		dbOrder.CreatedAt,
		dbOrder.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("order with uuid %s already exists", dbOrder.OrderUUID)
			}
		}
	}
	fmt.Printf("заказ создан: %d\n", orderM.OrderUUID)
	return nil
}
