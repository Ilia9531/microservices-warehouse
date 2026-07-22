package config

import "time"

// LoggerConfig определяет интерфейс для конфигурации логгера.
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// IamGRPCConfig определяет интерфейс для gRPC-конфигурации IAM.
type IamGRPCConfig interface {
	Host() string
	Port() string
	Address() string // Host:Port
}

// PostgresConfig определяет интерфейс для конфигурации PostgreSQL.
type PostgresConfig interface {
	DSN() string
	MigrationsDir() string
}

// RedisConfig определяет интерфейс для конфигурации Redis.
type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
	CacheTTL() time.Duration
}

// SessionConfig определяет интерфейс для конфигурации сессий.
type SessionConfig interface {
	TTLHour() time.Duration
}
