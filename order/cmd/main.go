package main

import (
	"github.com/Ilia9531/microservices-warehouse/order/internal/app"
	"github.com/Ilia9531/microservices-warehouse/order/internal/config"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/closer"
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
)

const configPath = "deploy/compose/order/.env"

func main() {
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	appCtx, appCancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	err = a.Run(appCtx)
	if err != nil {
		log.Fatalf("app run error: %v", err)
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
