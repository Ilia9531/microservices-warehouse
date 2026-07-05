package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"Jopa/assembly/internal/config"
	kafkaConverter "Jopa/assembly/internal/converter/kafka"
	"Jopa/assembly/internal/converter/kafka/decoder"
	"Jopa/assembly/internal/service"
	orderConsumer "Jopa/assembly/internal/service/consumer/order_consumer"
	orderProducer "Jopa/assembly/internal/service/producer/order_producer"
	"Jopa/platform/pkg/closer"
	wrappedKafka "Jopa/platform/pkg/kafka"
	wrappedKafkaConsumer "Jopa/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "Jopa/platform/pkg/kafka/producer"
	"Jopa/platform/pkg/logger"
	kafkaMiddleware "Jopa/platform/pkg/middleware/kafka"
)

type diContainer struct {
	// Kafka компоненты
	syncProducer           sarama.SyncProducer
	orderAssembledProducer wrappedKafka.Producer
	orderPaidConsumerGroup sarama.ConsumerGroup
	orderPaidConsumer      wrappedKafka.Consumer
	orderPaidDecoder       kafkaConverter.OrderPaidDecoder

	// Сервисы
	producerService service.OrderProducerService
	consumerService service.ConsumerService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		cfg := config.AppConfig().OrderAssembledProducer.Config()
		p, err := sarama.NewSyncProducer(config.AppConfig().Kafka.Brokers(), cfg)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer (ShipAssembled)", func(ctx context.Context) error {
			return p.Close()
		})
		d.syncProducer = p
	}
	return d.syncProducer
}

func (d *diContainer) OrderAssembledProducer() wrappedKafka.Producer {
	if d.orderAssembledProducer == nil {
		d.orderAssembledProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderAssembledProducer.Topic(),
			logger.Logger(),
		)
		logger.Info(context.Background(), "✅ ShipAssembled Kafka producer initialized")
	}
	return d.orderAssembledProducer
}

func (d *diContainer) OrderPaidConsumerGroup() sarama.ConsumerGroup {
	if d.orderPaidConsumerGroup == nil {
		cfg := config.AppConfig().OrderPaidConsumer.Config()
		cg, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			cfg,
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group (OrderPaid)", func(ctx context.Context) error {
			return cg.Close()
		})
		d.orderPaidConsumerGroup = cg
	}
	return d.orderPaidConsumerGroup
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidConsumerGroup(),
			[]string{config.AppConfig().OrderPaidConsumer.Topic()},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
		logger.Info(context.Background(), "✅ OrderPaid Kafka consumer initialized")
	}
	return d.orderPaidConsumer
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}
	return d.orderPaidDecoder
}

func (d *diContainer) ProducerService() service.OrderProducerService {
	if d.producerService == nil {
		d.producerService = orderProducer.NewService(d.OrderAssembledProducer())
		logger.Info(context.Background(), "✅ ProducerService initialized")
	}
	return d.producerService
}

func (d *diContainer) ConsumerService(ctx context.Context) service.ConsumerService {
	if d.consumerService == nil {
		d.consumerService = orderConsumer.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
			d.ProducerService(),
		)
		logger.Info(ctx, "✅ ConsumerService initialized")
	}
	return d.consumerService
}
