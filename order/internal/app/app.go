package app

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"Jopa/order/internal/config"
	"Jopa/platform/pkg/closer"
	"Jopa/platform/pkg/logger"
	orderv1 "Jopa/shared/pkg/openapi/order/v1"
)

// App — основной тип приложения.
type App struct {
	diContainer *diContainer
	httpServer  *http.Server
	listener    net.Listener // Явный слушатель
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
	// Канал для ошибок от компонентов
	errCh := make(chan error, 2)

	// Контекст для остановки всех горутин
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Консьюмер
	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- errors.Errorf("consumer crashed: %v", err)
		}
	}()

	go func() {
		err := a.runHTTPServer(ctx)
		if err != nil {
			errCh <- errors.Errorf("http server crashed: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		// Триггерим cancel, чтобы остановить второй компонент
		cancel()
		// Дождись завершения всех задач (если есть graceful shutdown внутри)
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
		a.runMigrations,
		a.initListener,   // Сначала открываем порт
		a.initHTTPServer, // Потом вешаем сервер на порт
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

func (a *App) initPostgres(ctx context.Context) error {
	a.diContainer.PostgresConn(ctx)
	return nil
}

func (a *App) runMigrations(ctx context.Context) error {
	a.diContainer.RunMigrations(ctx)
	return nil
}

// initListener создаёт TCP Listener и регистрирует его закрытие.
func (a *App) initListener(_ context.Context) error {
	cfg := config.AppConfig().Order

	// Открываем порт
	listener, err := net.Listen("tcp", cfg.Address())
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", cfg.Address(), err)
	}

	// Регистрируем закрытие Listener'а в closer
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

// initHTTPServer настраивает маршрутизацию и сервер.
func (a *App) initHTTPServer(ctx context.Context) error {
	// Создаём OpenAPI сервер на основе хендлера из DI
	apiServer, err := orderv1.NewServer(a.diContainer.APIHandler(ctx))
	if err != nil {
		return fmt.Errorf("failed to create openapi server: %w", err)
	}

	// Настраиваем роутер
	r := chi.NewRouter()
	r.Mount("/", apiServer)

	// Создаём HTTP сервер, привязывая его к нашему Listener'у
	a.httpServer = &http.Server{
		Handler:     r,
		ReadTimeout: config.AppConfig().Order.Timeout(),
		// Addr можно не указывать, если мы используем listener напрямую через Serve
	}

	// Регистрируем Graceful Shutdown сервера в closer
	// Shutdown ждёт завершения активных запросов
	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		err := a.httpServer.Shutdown(ctx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP OrderService server listening on %s", a.listener.Addr()))

	// Используем Serve с явным listener'ом вместо ListenAndServe
	err := a.httpServer.Serve(a.listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 UFORecorded Kafka consumer running")

	err := a.diContainer.OrderConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
