//go:build integration

package integration

import "time"

const (
	// projectName — имя проекта для Docker-контейнеров и сети
	projectName = "inventory-service"

	// partsCollectionName — имя коллекции MongoDB для деталей
	partsCollectionName = "parts"

	// Имена контейнеров (должны совпадать с конфигурацией в docker-compose)
	inventoryAppName    = "inventory-app"
	inventoryDockerfile = "deploy/docker/inventory/Dockerfile"

	// Ключи переменных окружения (из platform/pkg/testcontainers/path/constants.go)
	mongoImageNameKey = "MONGO_IMAGE_NAME"
	mongoHostKey      = "MONGO_HOST"
	mongoPortKey      = "MONGO_PORT"
	mongoDatabaseKey  = "MONGO_DATABASE"
	mongoUsernameKey  = "MONGO_INITDB_ROOT_USERNAME"
	mongoPasswordKey  = "MONGO_INITDB_ROOT_PASSWORD"
	mongoAuthDBKey    = "MONGO_AUTH_DB"
	grpcPortKey       = "GRPC_PORT"
	grpcHostKey       = "GRPC_HOST"

	// Значения по умолчанию
	grpcPortValue    = "50051"
	loggerLevelValue = "debug"

	startupTimeoutValue = 3 * time.Minute
)
