package config

import "github.com/IBM/sarama"

// KafkaConfig определяет контракт для конфигурации Kafka-брокеров.
type KafkaConfig interface {
	Brokers() []string
}

// LoggerConfig определяет контракт для конфигурации логгера.
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// OrderPaidConsumerConfig определяет контракт для конфигурации потребителя события OrderPaid.
type OrderPaidConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

// OrderAssembledProducerConfig определяет контракт для конфигурации продюсера события ShipAssembled.
// (Соответствует файлу env/order_assembled_producer.go, публикует ShipAssembled)
type OrderAssembledProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}
