package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type payGRPCEnvConfig struct {
	Host string `env:"PAYMENT_GRPC_HOST,required"`
	Port string `env:"PAYMENT_GRPC_PORT,required"`
}

type payGRPCConfig struct {
	raw payGRPCEnvConfig
}

func NewPaymentConfig() (*payGRPCConfig, error) {
	var raw payGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &payGRPCConfig{raw: raw}, nil
}

func (cfg *payGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
