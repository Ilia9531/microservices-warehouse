package service

import (
	"github.com/Ilia9531/microservices-warehouse/assembly/internal/model"
	"context"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderProducerService interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}
