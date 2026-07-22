package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/config/env"
)

// appConfig — глобальный синглтон конфигурации.
// Инициализируется один раз через Load() и доступен через AppConfig().
var appConfig *config

// config — внутренняя структура конфигурации.
// Поля экспортированы, так как реализуют интерфейсы, объявленные в interfaces.go.
type config struct {
	Logger   LoggerConfig
	GRPC     IamGRPCConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Session  SessionConfig
}

func Load(path ...string) error {
	// Загружаем .env файл, игнорируя ошибку "файл не найден"
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Инициализируем каждую подконфигурацию через пакет env
	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	grpcCfg, err := env.NewIamGRPCConfig()
	if err != nil {
		return err
	}

	pgCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	sessionCfg, err := env.NewSessionConfig()
	if err != nil {
		return err
	}

	// Сохраняем в глобальную переменную
	appConfig = &config{
		Logger:   loggerCfg,
		GRPC:     grpcCfg,
		Postgres: pgCfg,
		Redis:    redisCfg,
		Session:  sessionCfg,
	}

	return nil
}

// AppConfig возвращает глобальный экземпляр конфигурации.
// ⚠️ Важно: вызов до Load() вернёт nil.
// Рекомендуется использовать только после успешной инициализации.

func AppConfig() *config {
	return appConfig
}
