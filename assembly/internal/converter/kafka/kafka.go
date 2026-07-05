package kafka

import "Jopa/assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(payload []byte) (model.OrderPaidEvent, error)
}
