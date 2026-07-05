package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type OrderEnvConfig struct {
	Host    string        `env:"HTTP_HOST,required"`
	Port    string        `env:"HTTP_PORT,required"`
	Timeout time.Duration `env:"HTTP_READ_TIMEOUT,required"`
}

type OrderConfig struct {
	raw OrderEnvConfig
}

func NewOrderConfig() (*OrderConfig, error) {
	var raw OrderEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &OrderConfig{raw: raw}, nil
}

func (cfg *OrderConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *OrderConfig) Timeout() time.Duration {
	return cfg.raw.Timeout
}
