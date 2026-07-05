package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

//# Хост PostgreSQL-сервера (для внутренних подключений)
//POSTGRES_HOST=localhost
//
//# Внутренний порт PostgreSQL
//POSTGRES_PORT=5432
//
//# Внешний порт PostgreSQL (для подключения извне контейнера)
//EXTERNAL_POSTGRES_PORT=5432
//
//# Имя пользователя для подключения к PostgreSQL
//POSTGRES_USER=order_user
//
//# Пароль пользователя для подключения к PostgreSQL
//POSTGRES_PASSWORD=order_password
//
//# Название базы данных
//POSTGRES_DB=order
//
//# Режим подключения по SSL (например, disable, require)
//POSTGRES_SSL_MODE=disable
//
//# Путь к директории с миграциями
//MIGRATION_DIRECTORY=./order/migrations

type postgresEnvConfig struct {
	Host          string `env:"POSTGRES_HOST,required"`
	Port          string `env:"POSTGRES_PORT,required"`
	User          string `env:"POSTGRES_USER,required"`
	Password      string `env:"POSTGRES_PASSWORD,required"`
	Name          string `env:"POSTGRES_DB,required"`
	SSLMode       string `env:"POSTGRES_SSL_MODE,required"`
	MigrationsDir string `env:"MIGRATION_DIRECTORY,required"`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}
	return &postgresConfig{raw: raw}, nil
}

func (cfg *postgresConfig) DSN() string {
	//ORDER_POSTGRES_DSN=postgres://order-service-user:order-service-password@postgres:5432/order-service
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.raw.User,
		cfg.raw.Password,
		cfg.raw.Host,
		cfg.raw.Port,
		cfg.raw.Name,
		cfg.raw.SSLMode,
	)
}
func (cfg *postgresConfig) MigrationsDir() string {
	return cfg.raw.MigrationsDir
}
