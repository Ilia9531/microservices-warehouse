package config

import (
	"time"

	"github.com/IBM/sarama"
)

// LoggerConfig — интерфейс конфигурации логгера.
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// OrderHTTPConfig — интерфейс конфигурации HTTP-сервера Order.
type OrderHTTPConfig interface {
	Address() string // возвращает "host:port"
	Timeout() time.Duration
}

// PostgresConfig — интерфейс конфигурации PostgreSQL.
type PostgresConfig interface {
	DSN() string
	MigrationsDir() string
}

// InventoryGRPCConfig — интерфейс конфигурации клиента к Inventory gRPC сервису.
type InventoryGRPCConfig interface {
	Address() string // возвращает "host:port"
}

// PaymentGRPCConfig — интерфейс конфигурации клиента к Payment gRPC сервису.
type PaymentGRPCConfig interface {
	Address() string // возвращает "host:port"
}

// Интерфейс конфигурации кафки
type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderAssembledConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
