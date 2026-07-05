//go:build integration

package integration

import (
	"context"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go/wait"

	"go.uber.org/zap"

	"Jopa/platform/pkg/logger"
	"Jopa/platform/pkg/testcontainers"
	"Jopa/platform/pkg/testcontainers/app"
	"Jopa/platform/pkg/testcontainers/mongo"
	"Jopa/platform/pkg/testcontainers/network"
	"Jopa/platform/pkg/testcontainers/path"
)

// setupTestEnvironment подготавливает тестовое окружение: сеть, контейнеры и возвращает структуру с ресурсами
func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	log := logger.Logger()
	log.Info(ctx, "🚀 Подготовка тестового окружения...")

	// Шаг 1: Создаём общую Docker-сеть
	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "не удалось создать общую сеть", zap.Error(err))
	}
	log.Info(ctx, "✅ Сеть успешно создана", zap.String("name", generatedNetwork.Name()))

	// Получаем переменные окружения для MongoDB с проверкой на наличие
	mongoImageName := getEnvWithLogging(ctx, testcontainers.MongoImageNameKey)
	mongoUsername := getEnvWithLogging(ctx, testcontainers.MongoUsernameKey)
	mongoPassword := getEnvWithLogging(ctx, testcontainers.MongoPasswordKey)
	mongoDatabase := getEnvWithLogging(ctx, testcontainers.MongoDatabaseKey)
	mongoAuthDB := getEnvWithLogging(ctx, testcontainers.MongoAuthDBKey)

	// Получаем порт gRPC для waitStrategy
	grpcPort := getEnvWithLogging(ctx, grpcPortKey)
	if grpcPort == "" {
		grpcPort = grpcPortValue // fallback из constants.go
	}

	// Шаг 2: Запускаем контейнер с MongoDB
	log.Info(ctx, "🐳 Запускаем MongoDB контейнер...")
	generatedMongo, err := mongo.NewContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithImageName(mongoImageName),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithAuthDB(mongoAuthDB),
		mongo.WithLogger(log),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер MongoDB", zap.Error(err))
	}
	log.Info(ctx, "✅ Контейнер MongoDB успешно запущен",
		zap.String("host", generatedMongo.Config().Host),
		zap.String("port", generatedMongo.Config().Port),
	)

	// Шаг 3: Запускаем контейнер с приложением Inventory
	projectRoot := path.GetProjectRoot()

	appEnv := map[string]string{
		// Переопределяем URI MongoDB для подключения из контейнера приложения
		testcontainers.MongoHostKey: generatedMongo.Config().ContainerName,
	}

	// Создаём настраиваемую стратегию ожидания с увеличенным таймаутом
	// Ждём, пока порт 50051 станет доступен для gRPC-соединений
	//waitStrategy := wait.ForListeningPort(grpcPort + "/tcp").
	//	WithStartupTimeout(startupTimeoutValue)
	waitStrategy := wait.ForLog("🚀 gRPC InventoryService server listening on").WithStartupTimeout(110 * time.Second)

	log.Info(ctx, "📦 Собираем и запускаем Inventory контейнер...",
		zap.String("dockerfile", inventoryDockerfile),
		zap.String("grpc_port", grpcPort),
	)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(inventoryAppName),
		app.WithPort(grpcPort),
		app.WithDockerfile(projectRoot, inventoryDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(log),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{
			Network: generatedNetwork,
			Mongo:   generatedMongo,
		})
		logger.Fatal(ctx, "не удалось запустить контейнер приложения", zap.Error(err))
	}
	log.Info(ctx, "✅ Контейнер приложения успешно запущен",
		zap.String("address", appContainer.Address()),
	)

	log.Info(ctx, "🎉 Тестовое окружение готово")

	return &TestEnvironment{
		Network: generatedNetwork,
		Mongo:   generatedMongo,
		App:     appContainer,
	}
}

// getEnvWithLogging возвращает значение переменной окружения с логированием
func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Warn(ctx, "Переменная окружения не установлена", zap.String("key", key))
	}
	return value
}
