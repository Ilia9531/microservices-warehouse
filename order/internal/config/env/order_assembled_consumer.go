package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

// orderAssembledConsumerEnvConfig — сырая конфигурация из переменных окружения.
// Связана с топиком order.assembled, который публикует AssemblyService.
type orderAssembledConsumerEnvConfig struct {
	Topic   string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_ASSEMBLED_CONSUMER_GROUP_ID,required"`
}

// orderAssembledConsumerConfig — публичный конфиг, инкапсулирующий raw-данные.
type orderAssembledConsumerConfig struct {
	raw orderAssembledConsumerEnvConfig
}

// NewOrderAssembledConsumerConfig парсит переменные окружения и возвращает готовый конфиг.
func NewOrderAssembledConsumerConfig() (*orderAssembledConsumerConfig, error) {
	var raw orderAssembledConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderAssembledConsumerConfig{raw: raw}, nil
}

// Topic возвращает имя Kafka-топика для подписки на события ShipAssembled.
func (cfg *orderAssembledConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

// GroupID возвращает ID consumer group для координации между инстансами OrderService.
func (cfg *orderAssembledConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

// Config возвращает настроенный экземпляр sarama.Config для consumer.
func (cfg *orderAssembledConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
