package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

// orderPaidConsumerEnvConfig — сырая конфигурация из переменных окружения.
// Связана с топиком order.assembled, который публикует AssemblyService.
type orderPaidConsumerEnvConfig struct {
	Topic   string `env:"ORDER_PAID_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_PAID_CONSUMER_GROUP_ID,required"`
}

// orderPaidConsumerConfig — публичный конфиг, инкапсулирующий raw-данные.
type orderPaidConsumerConfig struct {
	raw orderPaidConsumerEnvConfig
}

// NewOrderPaidConsumerConfig парсит переменные окружения и возвращает готовый конфиг.
func NewOrderPaidConsumerConfig() (*orderPaidConsumerConfig, error) {
	var raw orderPaidConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPaidConsumerConfig{raw: raw}, nil
}

// Topic возвращает имя Kafka-топика для подписки на события ShipAssembled.
func (cfg *orderPaidConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

// GroupID возвращает ID consumer group для координации между инстансами OrderService.
func (cfg *orderPaidConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

// Config возвращает настроенный экземпляр sarama.Config для consumer.
func (cfg *orderPaidConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
