package kafka

import "github.com/Ilia9531/microservices-warehouse/assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(payload []byte) (model.OrderPaidEvent, error)
}
