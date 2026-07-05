package kafka

import "Jopa/order/internal/model"

type ShipAssembledDecoder interface {
	Decode(payload []byte) (model.ShipAssembledEvent, error)
}
