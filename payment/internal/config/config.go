package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Ilia9531/microservices-warehouse/payment/internal/config/env"
)

// appConfig — глобальный синглтон конфигурации.
// Инициализируется один раз через Load() и доступен через AppConfig().
var appConfig *config

// config — внутренняя структура конфигурации.
// Поля экспортированы, так как реализуют интерфейсы, объявленные в interfaces.go.
type config struct {
	Logger LoggerConfig
	GRPC   PaymentGRPCConfig
}

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

	grpcCfg, err := env.NewPaymentConfig()
	if err != nil {
		return err
	}

	// Сохраняем в глобальную переменную
	appConfig = &config{
		Logger: loggerCfg,
		GRPC:   grpcCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
