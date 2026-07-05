package mongo

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	mytestcontainers "github.com/Ilia9531/microservices-warehouse/platform/pkg/testcontainers"
)

const (
	mongoPort           = mytestcontainers.MongoPort
	mongoStartupTimeout = 100 * time.Second

	mongoEnvUsernameKey = mytestcontainers.MongoUsernameKey
	mongoEnvPasswordKey = mytestcontainers.MongoPasswordKey //nolint:gosec
)

type Container struct {
	container testcontainers.Container
	client    *mongo.Client
	cfg       *Config
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := buildConfig(opts...)

	container, err := startMongoContainer(ctx, cfg)
	if err != nil {
		return nil, err
	}

	success := false
	defer func() {
		if !success {
			if err = container.Terminate(ctx); err != nil {
				cfg.Logger.Error(ctx, "failed to terminate mongo container", zap.Error(err))
			}
		}
	}()

	cfg.Host, cfg.Port, err = getContainerHostPort(ctx, container)
	if err != nil {
		cfg.Logger.Error(ctx, "MappedPort failed, switching to internal Docker network mode",
			zap.String("container", cfg.ContainerName),
			zap.String("internal_port", mongoPort),
		)
		return nil, err
		// cfg.Host = cfg.ContainerName // Имя контейнера резолвится внутри Docker-сети как DNS-хост
		// cfg.Port = mongoPort
	}

	uri := buildMongoURI(cfg)
	client, err := connectMongoClient(ctx, uri)
	if err != nil {
		return nil, err
	}

	cfg.Logger.Info(ctx, "Mongo container started, ВОТ ПУТЬ", zap.String("uri", uri))
	success = true

	return &Container{
		container: container,
		client:    client,
		cfg:       cfg,
	}, nil
}

func (c *Container) Client() *mongo.Client {
	return c.client
}

func (c *Container) Config() *Config {
	return c.cfg
}

func (c *Container) Terminate(ctx context.Context) error {
	if err := c.client.Disconnect(ctx); err != nil {
		c.cfg.Logger.Error(ctx, "failed to disconnect mongo client", zap.Error(err))
	}

	if err := c.container.Terminate(ctx); err != nil {
		c.cfg.Logger.Error(ctx, "failed to terminate mongo container", zap.Error(err))
	}

	c.cfg.Logger.Info(ctx, "Mongo container terminated")

	return nil
}
