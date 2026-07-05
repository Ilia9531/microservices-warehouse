package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type invGRPCEnvConfig struct {
	Host string `env:"INVENTORY_GRPC_HOST,required"`
	Port string `env:"INVENTORY_GRPC_PORT,required"`
}

type invGRPCConfig struct {
	raw invGRPCEnvConfig
}

func NewInventoryConfig() (*invGRPCConfig, error) {
	var raw invGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &invGRPCConfig{raw: raw}, nil
}

func (cfg *invGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
