package kafka

import "github.com/Ilia9531/microservices-warehouse/order/internal/model"

type ShipAssembledDecoder interface {
	Decode(payload []byte) (model.ShipAssembledEvent, error)
}
