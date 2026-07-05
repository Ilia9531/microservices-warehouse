package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

// orderPaidProducerEnvConfig — маппинг переменных окружения на структуру.
// Тег required гарантирует, что сервис не запустится без имени топика.
type orderPaidProducerEnvConfig struct {
	TopicName string `env:"ORDER_PAID_TOPIC_NAME,required"`
}

// orderPaidProducerConfig — публичный конфиг, скрывающий детали парсинга.
type orderPaidProducerConfig struct {
	raw orderPaidProducerEnvConfig
}

// NewOrderPaidProducerConfig читает переменные окружения и возвращает готовый конфиг.
func NewOrderPaidProducerConfig() (*orderPaidProducerConfig, error) {
	var raw orderPaidProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPaidProducerConfig{raw: raw}, nil
}

// Topic возвращает имя Kafka-топика для публикации события OrderPaid.
func (cfg *orderPaidProducerConfig) Topic() string {
	return cfg.raw.TopicName
}

// Config возвращает настройки Sarama-продюсера.
// Return.Successes = true позволяет синхронно ждать подтверждения от брокера.
func (cfg *orderPaidProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
