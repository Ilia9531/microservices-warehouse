package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/Ilia9531/microservices-warehouse/inventory/internal/app"
	"github.com/Ilia9531/microservices-warehouse/inventory/internal/config"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/closer"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/logger"
)

// если запуск локально: deploy/compose/inventory/
// если запуск из докер контейнера: .env
const configPath = "deploy/compose/inventory/.env"

func main() {
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		fmt.Printf("🔍 godotenv error type: %T\n", err)
		fmt.Printf("🔍 godotenv error msg: %v\n", err)

		panic(fmt.Errorf("failed to load .env file: %w", err))
	}
	appCtx, appCancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось создать приложение", zap.Error(err))
		return
	}
	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Ошибка при работе приложения", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
