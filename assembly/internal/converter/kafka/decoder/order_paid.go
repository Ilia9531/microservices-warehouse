package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/Ilia9531/microservices-warehouse/assembly/internal/model"
	eventsv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderPaidDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(payload []byte) (model.OrderPaidEvent, error) {
	var pb eventsv1.OrderPaidEvent
	if err := proto.Unmarshal(payload, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("unmarshal order paid event: %w", err)
	}
	return model.OrderPaidEvent{
		EventUUID:       pb.EventUuid,
		OrderUUID:       pb.OrderUuid,
		UserUUID:        pb.UserUuid,
		PaymentMethod:   pb.PaymentMethod,
		TransactionUUID: pb.TransactionUuid,
	}, nil
}
