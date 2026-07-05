package model

import "errors"

var (
	ErrInvalidOrderUUID     = errors.New("invalid order uuid")
	ErrInvalidUserUUID      = errors.New("invalid user uuid")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
)
