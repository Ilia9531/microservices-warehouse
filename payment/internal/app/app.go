// payment/internal/app/app.go

package app

import (
	"context"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"Jopa/payment/internal/config"
	"Jopa/platform/pkg/closer"
	"Jopa/platform/pkg/grpc/health"
	"Jopa/platform/pkg/logger"
	payV1 "Jopa/shared/pkg/proto/payment/v1"
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

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run запускает приложение.
func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

// initDeps последовательно инициализирует все зависимости приложения.
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initListener,   // 1. Открываем порт
		a.initGRPCServer, // 2. Настраиваем сервер
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
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJSON(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

// initListener открывает TCP порт и регистрирует закрытие сокета.
func (a *App) initListener(_ context.Context) error {
	cfg := config.AppConfig().GRPC

	listener, err := net.Listen("tcp", cfg.Address())
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", cfg.Address(), err)
	}

	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		err := listener.Close()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			return err
		}
		return nil
	})

	a.listener = listener
	return nil
}

// initGRPCServer настраивает gRPC сервер, регистрирует сервисы и health check.
func (a *App) initGRPCServer(ctx context.Context) error {
	logger.Info(ctx, "🔧 [1/5] Creating gRPC server...")
	// Создаём сервер с insecure credentials (для локальной разработки)
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	logger.Info(ctx, "🔧 [2/5] Registering closer...")
	// Регистрируем Graceful Shutdown
	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})
	logger.Info(ctx, "🔧 [3/5] Registering reflection...")
	// Reflection нужен для grpcurl (тестирование через task test-api)
	reflection.Register(a.grpcServer)
	logger.Info(ctx, "🔧 [4/5] Getting API handler...")
	// Health Check — обязателен для production readiness
	health.RegisterService(a.grpcServer)
	logger.Info(ctx, "🔧 [5/5] Registering PaymentService...")
	// Регистрируем наш бизнес-сервис
	payV1.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.APIHandler(ctx))
	logger.Info(ctx, "✅ PaymentService successfully registered")
	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf(" Payment gRPC server listening on %s", a.listener.Addr()))

	// Serve блокирует выполнение до закрытия сервера
	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}
