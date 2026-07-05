package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"Jopa/order/internal/model"
	eventsv1 "Jopa/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewShipAssembledDecoder() *decoder {
	return &decoder{}
}

// Decode преобразует сырой payload Kafka-сообщения в доменное событие.
// Возвращает ошибку, если десериализация protobuf завершилась неудачей.
func (d *decoder) Decode(payload []byte) (model.ShipAssembledEvent, error) {
	var pbEvent eventsv1.ShipAssembledEvent
	if err := proto.Unmarshal(payload, &pbEvent); err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("unmarshal ship assembled event: %w", err)
	}

	return model.ShipAssembledEvent{
		EventUUID:    pbEvent.EventUuid,
		OrderUUID:    pbEvent.OrderUuid,
		UserUUID:     pbEvent.UserUuid,
		BuildTimeSec: pbEvent.BuildTimeSec,
	}, nil
}
