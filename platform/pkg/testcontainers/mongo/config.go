package mongo

import (
	"context"

	//"github.com/docker/go-connections/nat"
	"github.com/docker/docker/api/types/container"
	"go.uber.org/zap"

	//"github.com/moby/moby/api/types/network"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/logger"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Config struct {
	NetworkName   string
	ContainerName string
	ImageName     string
	Database      string
	Username      string
	Password      string
	AuthDB        string
	Logger        Logger

	Host string
	Port string
}

func buildConfig(opts ...Option) *Config {
	cfg := &Config{
		NetworkName:   "inventory-service",
		ContainerName: "mongo",
		ImageName:     "mongo:7.0.5",
		Database:      "inventory",
		Username:      "inventory_admin",
		Password:      "inventory_secret",
		AuthDB:        "admin",
		Logger:        &logger.NoopLogger{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func defaultHostConfig() func(hc *container.HostConfig) {
	return func(hc *container.HostConfig) {
		hc.AutoRemove = false
	}
}
