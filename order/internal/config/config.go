package config

import (
	"os"

	"github.com/joho/godotenv"

	"Jopa/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	Order                  OrderHTTPConfig
	InvGRPC                InventoryGRPCConfig
	PayGRPC                PaymentGRPCConfig
	Postgres               PostgresConfig
	Kafka                  KafkaConfig
	OrderPaidProducer      OrderPaidProducerConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Инициализируем каждую подконфигурацию через пакет env
	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	orderCfg, err := env.NewOrderConfig()
	if err != nil {
		return err
	}

	InvGRPCCfg, err := env.NewInventoryConfig()
	if err != nil {
		return err
	}

	PayGRPCCfg, err := env.NewPaymentConfig()
	if err != nil {
		return err
	}

	PostgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	OrderPaidProducerCfg, err := env.NewOrderPaidProducerConfig()
	if err != nil {
		return err
	}

	OrderAssembledConsumerCfg, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		Order:                  orderCfg,
		InvGRPC:                InvGRPCCfg,
		PayGRPC:                PayGRPCCfg,
		Postgres:               PostgresCfg,
		Kafka:                  kafkaCfg,
		OrderPaidProducer:      OrderPaidProducerCfg,
		OrderAssembledConsumer: OrderAssembledConsumerCfg,
	}
	return nil
}

func AppConfig() *config {
	return appConfig
}
