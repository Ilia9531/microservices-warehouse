package app

import (
	"context"
	"fmt"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	authv1 "github.com/Ilia9531/microservices-warehouse/iam/internal/api/auth/v1"
	userv1 "github.com/Ilia9531/microservices-warehouse/iam/internal/api/user/v1"
	"github.com/Ilia9531/microservices-warehouse/iam/internal/config"
	"github.com/Ilia9531/microservices-warehouse/iam/internal/repository"
	sessionRepo "github.com/Ilia9531/microservices-warehouse/iam/internal/repository/session"
	userRepo "github.com/Ilia9531/microservices-warehouse/iam/internal/repository/user"
	Svc "github.com/Ilia9531/microservices-warehouse/iam/internal/service"
	authSvc "github.com/Ilia9531/microservices-warehouse/iam/internal/service/auth"
	userSvc "github.com/Ilia9531/microservices-warehouse/iam/internal/service/user"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/cache"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/cache/redis"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/closer"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/logger"
	pgMigrator "github.com/Ilia9531/microservices-warehouse/platform/pkg/migrator"
)

type diContainer struct {
	pgConn     *pgx.Conn
	redisCache cache.RedisClient

	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository

	userService Svc.UserService
	authService Svc.AuthService

	userAPI *userv1.Api
	authAPI *authv1.Api
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PostgresConn(ctx context.Context) *pgx.Conn {
	if d.pgConn == nil {
		cfg := config.AppConfig().Postgres
		conn, err := pgx.Connect(ctx, cfg.DSN())
		if err != nil {
			panic(fmt.Sprintf("❌ Failed to connect to PostgreSQL: %v", err))
		}
		if err = conn.Ping(ctx); err != nil {
			panic(fmt.Sprintf("❌ PostgreSQL ping failed: %v", err))
		}
		closer.AddNamed("Postgres client", func(ctx context.Context) error {
			return conn.Close(ctx)
		})
		logger.Info(ctx, "✅ Connected to PostgreSQL")
		d.pgConn = conn
	}
	return d.pgConn
}

func (d *diContainer) RunMigrations(ctx context.Context) {
	cfg := config.AppConfig().Postgres
	db := stdlib.OpenDB(*d.PostgresConn(ctx).Config().Copy())
	defer func() { _ = db.Close() }()

	mig := pgMigrator.NewMigrator(db, cfg.MigrationsDir())
	if err := mig.Up(); err != nil {
		panic(fmt.Sprintf("failed to run migrations: %v", err))
	}
	logger.Info(ctx, "✅ Migrations completed")
}

func (d *diContainer) RedisCache(ctx context.Context) cache.RedisClient {
	if d.redisCache == nil {
		cfg := config.AppConfig().Redis
		pool := &redigo.Pool{
			MaxIdle:     cfg.MaxIdle(),
			IdleTimeout: cfg.IdleTimeout(),
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", cfg.Address())
			},
			TestOnBorrow: func(c redigo.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}
		d.redisCache = redis.NewClient(pool, logger.Logger(), cfg.ConnectionTimeout())
		if err := d.redisCache.Ping(ctx); err != nil {
			panic(fmt.Sprintf("❌ Failed to connect to Redis: %v", err))
		}
		closer.AddNamed("Redis client", func(ctx context.Context) error {
			// redigo.Pool не имеет Close, но клиент закрывает соединения сам
			return nil
		})
		logger.Info(ctx, "✅ Connected to Redis")
	}
	return d.redisCache
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepo == nil {
		d.userRepo = userRepo.NewRepository(d.PostgresConn(ctx))
		logger.Info(ctx, "✅ UserRepository initialized")
	}
	return d.userRepo
}

func (d *diContainer) SessionRepository(ctx context.Context) repository.SessionRepository {
	if d.sessionRepo == nil {
		cfg := config.AppConfig().Session
		d.sessionRepo = sessionRepo.NewRepository(
			d.RedisCache(ctx),
			cfg.TTLHour(),
		)
		logger.Info(ctx, "✅ SessionRepository initialized")
	}
	return d.sessionRepo
}

func (d *diContainer) UserService(ctx context.Context) Svc.UserService {
	if d.userService == nil {
		d.userService = userSvc.NewService(d.UserRepository(ctx))
		logger.Info(ctx, "✅ UserService initialized")
	}
	return d.userService
}

func (d *diContainer) AuthService(ctx context.Context) Svc.AuthService {
	if d.authService == nil {
		cfg := config.AppConfig().Session
		d.authService = authSvc.NewService(
			d.UserRepository(ctx),
			d.SessionRepository(ctx),
			cfg.TTLHour(),
		)
		logger.Info(ctx, "✅ AuthService initialized")
	}
	return d.authService
}

func (d *diContainer) UserAPI(ctx context.Context) *userv1.Api {
	if d.userAPI == nil {
		d.userAPI = userv1.NewAPI(d.UserService(ctx))
		logger.Info(ctx, "✅ UserAPI initialized")
	}
	return d.userAPI
}

func (d *diContainer) AuthAPI(ctx context.Context) *authv1.Api {
	if d.authAPI == nil {
		d.authAPI = authv1.NewAPI(d.AuthService(ctx))
		logger.Info(ctx, "✅ AuthAPI initialized")
	}
	return d.authAPI
}
