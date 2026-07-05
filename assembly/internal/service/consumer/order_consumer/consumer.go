package order_consumer

import (
	serviceKafka "Jopa/assembly/internal/service"
	"context"

	"go.uber.org/zap"

	kafakDec "Jopa/assembly/internal/converter/kafka"
	"Jopa/platform/pkg/kafka"
	"Jopa/platform/pkg/logger"
)

// service реализует интерфейс ConsumerService для обработки входящих событий OrderPaid.
type service struct {
	consumer         kafka.Consumer
	orderPaidDecoder kafakDec.OrderPaidDecoder
	producerService  serviceKafka.OrderProducerService
}

// NewService создаёт новый экземпляр сервиса консьюмера.
func NewService(consumer kafka.Consumer, decoder kafakDec.OrderPaidDecoder,
	producer serviceKafka.OrderProducerService) *service {
	return &service{
		consumer:         consumer,
		orderPaidDecoder: decoder,
		producerService:  producer,
	}
}

// RunConsumer запускает цикл потребления сообщений из Kafka.
// Блокирует выполнение до отмены контекста или ошибки.
func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order paid consumer service")

	err := s.consumer.Consume(ctx, s.OrderPaidHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
