package config

// LoggerConfig — интерфейс конфигурации логгера.
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type PaymentGRPCConfig interface {
	Address() string // возвращает "host:port"
}
