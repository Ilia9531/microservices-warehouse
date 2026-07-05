package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/Ilia9531/microservices-warehouse/order/internal/converter/kafka"
	repo "github.com/Ilia9531/microservices-warehouse/order/internal/repository"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/kafka"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/logger"
)

// Service реализует интерфейс обработки Kafka-событий для OrderService.
type service struct {
	consumer             kafka.Consumer
	shipAssembledDecoder kafkaConverter.ShipAssembledDecoder
	orderRepo            repo.OrderRepository
}

// NewService создаёт новый экземпляр сервиса консьюмера.
func NewService(consumer kafka.Consumer, decoder kafkaConverter.ShipAssembledDecoder,
	orderRepo repo.OrderRepository,
) *service {
	return &service{
		consumer:             consumer,
		shipAssembledDecoder: decoder,
		orderRepo:            orderRepo,
	}
}

// RunConsumer запускает цикл потребления сообщений из Kafka.
// Блокирует выполнение до отмены контекста или ошибки.
func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order assembled consumer service")

	err := s.consumer.Consume(ctx, s.ShipAssembledHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.assembled topic error", zap.Error(err))
		return err
	}

	return nil
}
