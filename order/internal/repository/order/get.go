package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/Ilia9531/microservices-warehouse/order/internal/model"
	"github.com/Ilia9531/microservices-warehouse/order/internal/repository/converter"
	repoModel "github.com/Ilia9531/microservices-warehouse/order/internal/repository/model"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) Get(ctx context.Context, uuid string) (*model.Order, error) {
	if uuid == "" {
		return nil, fmt.Errorf("uuid is required")
	}
	query := ` SELECT * FROM orders WHERE order_uuid = $1`

	row := r.db.QueryRow(ctx, query, uuid)

	var o repoModel.Order
	err := row.Scan(
		&o.OrderUUID,
		&o.UserUUID,
		&o.PartUUIDs,
		&o.TotalPrice,
		&o.Status,
		&o.TransactionUUID,
		&o.PaymentMethod,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("order not found: %s", uuid)
		}
		return nil, model.ErrNotFound
	}
	return converter.ToDomain(&o), nil
}
