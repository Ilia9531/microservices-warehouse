package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

// orderAssembledProducerEnvConfig — маппинг переменных окружения на структуру.
// Тег required гарантирует, что сервис не запустится без имени топика.
type orderAssembledProducerEnvConfig struct {
	TopicName string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
}

// orderAssembledProducerConfig — публичный конфиг, скрывающий детали парсинга.
type orderAssembledProducerConfig struct {
	raw orderAssembledProducerEnvConfig
}

// NewOrderAssembledProducerConfig читает переменные окружения и возвращает готовый конфиг.
func NewOrderAssembledProducerConfig() (*orderAssembledProducerConfig, error) {
	var raw orderAssembledProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderAssembledProducerConfig{raw: raw}, nil
}

// Topic возвращает имя Kafka-топика для публикации события OrderAssembled.
func (cfg *orderAssembledProducerConfig) Topic() string {
	return cfg.raw.TopicName
}

// Config возвращает настройки Sarama-продюсера.
// Return.Successes = true позволяет синхронно ждать подтверждения от брокера.
func (cfg *orderAssembledProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
