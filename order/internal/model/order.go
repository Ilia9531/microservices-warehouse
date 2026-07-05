package model

import "time"

type Order struct {
	OrderUUID       string
	UserUUID        string
	PartUUIDs       []string
	TotalPrice      float64
	TransactionUUID *string     // pointer, т.к. опционально
	PaymentMethod   *string     // pointer, т.к. опционально
	Status          OrderStatus // PENDING_PAYMENT, PAID, CANCELLED, ASSEMBLED
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderStatus string

const (
	StatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	StatusPaid           OrderStatus = "PAID"
	StatusCancelled      OrderStatus = "CANCELLED"
	StatusAssembled      OrderStatus = "ASSEMBLED"
)
