//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/Ilia9531/microservices-warehouse/platform/pkg/logger"
)

const testsTimeout = 14 * time.Minute

var (
	env *TestEnvironment

	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Inventory Service Integration Test Suite")
}

var _ = BeforeSuite(func() {
	// Инициализируем логгер для тестов
	err := logger.Init(loggerLevelValue, true)
	if err != nil {
		panic(fmt.Sprintf("не удалось инициализировать логгер: %v", err))
	}

	// Создаем контекст с общим таймаутом на весь набор тестов
	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

	// Загружаем .env файл из deploy/compose/inventory/.env
	envPath := filepath.Join("..", "..", "..", "deploy", "compose", "inventory", ".env")

	envVars, err := godotenv.Read(envPath)
	if err != nil {
		logger.Fatal(suiteCtx, "Не удалось загрузить .env файл", zap.Error(err))
	}
	// Устанавливаем переменные в окружение процесса, чтобы setup.go их увидел
	for key, value := range envVars {
		_ = os.Setenv(key, value)
	}

	logger.Info(suiteCtx, "🚀 Запуск тестового окружения...")
	env = setupTestEnvironment(suiteCtx)
})

var _ = AfterSuite(func() {
	logger.Info(context.Background(), "🏁 Завершение набора тестов")
	if env != nil {
		teardownTestEnvironment(suiteCtx, env)
	}
	suiteCancel()
})
