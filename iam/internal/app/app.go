package app

import (
	"context"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/config"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/closer"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/grpc/health"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/logger"
	authv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/auth/v1"
	userv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/user/v1"
)

// App — основной тип приложения.
type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

// New создаёт новый экземпляр приложения и инициализирует все зависимости.
func New(ctx context.Context) (*App, error) {
	a := &App{}
	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}
	return a, nil
}

// Run запускает приложение.
func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := a.runGRPCServer(ctx); err != nil {
			errCh <- errors.Errorf("gRPC server crashed: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		cancel()
		<-ctx.Done()
		return err
	}
	return nil
}

// initDeps последовательно инициализирует все зависимости приложения.
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initPostgres,
		a.initRedis,
		a.runMigrations,
		a.initListener,
		a.initGRPCServer,
	}
	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	cfg := config.AppConfig().Logger
	return logger.Init(cfg.Level(), cfg.AsJSON())
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initPostgres(ctx context.Context) error {
	a.diContainer.PostgresConn(ctx)
	return nil
}

func (a *App) initRedis(ctx context.Context) error {
	a.diContainer.RedisCache(ctx)
	return nil
}

func (a *App) runMigrations(ctx context.Context) error {
	a.diContainer.RunMigrations(ctx)
	return nil
}

func (a *App) initListener(_ context.Context) error {
	cfg := config.AppConfig().GRPC
	listener, err := net.Listen("tcp", cfg.Address())
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", cfg.Address(), err)
	}
	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		if err := listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			return err
		}
		return nil
	})
	a.listener = listener
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer()

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})
	// Включаем reflection для grpcurl и отладки
	reflection.Register(a.grpcServer)

	health.RegisterService(a.grpcServer)

	// Регистрируем сервисы
	userv1.RegisterUserServiceServer(a.grpcServer, a.diContainer.UserAPI(ctx))
	authv1.RegisterAuthServiceServer(a.grpcServer, a.diContainer.AuthAPI(ctx))

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC IAM server listening on %s", a.listener.Addr()))
	if err := a.grpcServer.Serve(a.listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return err
	}
	return nil
}
