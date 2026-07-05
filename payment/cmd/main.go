package main

import (
	"github.com/Ilia9531/microservices-warehouse/payment/internal/app"
	"github.com/Ilia9531/microservices-warehouse/payment/internal/config"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/closer"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/logger"
	"context"
	"time"

	"go.uber.org/zap"

	"fmt"

	"os/signal"
	"syscall"
)

const configPath = "deploy/compose/payment/.env"

func main() {
	ctx := context.Background()

	err := config.Load(configPath)

	if err != nil {

		panic(fmt.Errorf("failed to load .env file: %v\n", err))
		return
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
