package config

import (
	"os"

	"github.com/joho/godotenv"

	"Jopa/inventory/internal/config/env"
)

// appConfig — глобальный синглтон конфигурации.
// Инициализируется один раз через Load() и доступен через AppConfig().
var appConfig *config

// config — внутренняя структура конфигурации.
// Поля экспортированы, так как реализуют интерфейсы, объявленные в interfaces.go.
type config struct {
	Logger LoggerConfig
	GRPC   InventoryGRPCConfig
	Mongo  MongoConfig
}

// Load загружает переменные окружения из .env-файлов и инициализирует глобальную конфигурацию.
//
// Параметры:
//   - path: список путей к .env-файлам (опционально).
//     Если не передан, godotenv попытается загрузить .env из текущей директории.
//
// Возвращает ошибку, если:
//   - godotenv не смог прочитать файл (кроме ошибки "файл не существует")
//   - одна из подконфигураций не прошла валидацию при создании
//
// Пример использования в main.go:
//
//	if err := config.Load("../../deploy/env/inventory.env"); err != nil {
//	    log.Fatalf("failed to load config: %v", err)
//	}
func Load(path ...string) error {
	// Загружаем .env файл, игнорируя ошибку "файл не найден"
	// Это позволяет запускать приложение без .env, используя только os.Getenv
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Инициализируем каждую подконфигурацию через пакет env
	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	grpcCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	// Сохраняем в глобальную переменную
	appConfig = &config{
		Logger: loggerCfg,
		GRPC:   grpcCfg,
		Mongo:  mongoCfg,
	}

	return nil
}

// AppConfig возвращает глобальный экземпляр конфигурации.
//
// ⚠️ Важно: вызов до Load() вернёт nil.
// Рекомендуется использовать только после успешной инициализации.
//
// Пример:
//
//	cfg := config.AppConfig()
//	logger, _ := zap.NewProductionConfig().Build(
//	    zap.WithLoggerConfig(cfg.Logger),
//	)
func AppConfig() *config {
	return appConfig
}
