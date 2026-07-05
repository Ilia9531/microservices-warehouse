package model

import "errors"

var (
	ErrPartNotFound = errors.New("part not found")
	ErrValidation   = errors.New("error validation")
)
