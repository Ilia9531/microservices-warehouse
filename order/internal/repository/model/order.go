package model

import "time"

type Order struct {
	OrderUUID       string    `db:"order_uuid"`
	UserUUID        string    `db:"user_uuid"`
	PartUUIDs       []string  `db:"part_uuids"`
	TotalPrice      float64   `db:"total_price"`
	TransactionUUID *string   `db:"transaction_uuid"` // pointer, т.к. опционально
	PaymentMethod   *string   `db:"payment_method"`   // pointer, т.к. опционально
	Status          string    `db:"status"`           // PENDING_PAYMENT, PAID, CANCELLED
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}
