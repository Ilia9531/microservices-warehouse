package config

// Logger LoggerConfig
// GRPC   InventoryGRPCConfig
// Mongo  MongoConfig
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// InventoryGRPCConfig — интерфейс конфигурации gRPC-сервера Inventory.
type InventoryGRPCConfig interface {
	Address() string // возвращает "host:port"
}

// MongoConfig — интерфейс конфигурации MongoDB.
type MongoConfig interface {
	URI() string
	DatabaseName() string
}

type IamGRPCConfig interface {
	Address() string
}
