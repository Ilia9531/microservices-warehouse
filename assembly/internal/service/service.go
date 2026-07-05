package service

import (
	"context"

	"github.com/Ilia9531/microservices-warehouse/assembly/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderProducerService interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}
