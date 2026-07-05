// payment/internal/app/di.go

package app

import (
	api "Jopa/payment/internal/api/payment/v1"
	"context"

	"Jopa/payment/internal/service"
	paymentSvc "Jopa/payment/internal/service/payment"
	"Jopa/platform/pkg/logger"
	payV1 "Jopa/shared/pkg/proto/payment/v1"
)

// diContainer — контейнер зависимостей для Payment-сервиса.
type diContainer struct {
	paymentService service.PaymentService
	apiHandler     payV1.PaymentServiceServer
}

// NewDiContainer создаёт новый экземпляр DI-контейнера.
func NewDiContainer() *diContainer {
	return &diContainer{}
}

// PaymentService инициализирует бизнес-логику сервиса оплаты.
func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentSvc.NewService()
		logger.Info(ctx, "✅ Payment service initialized")
	}
	return d.paymentService
}

// APIHandler инициализирует gRPC-хендлер.
func (d *diContainer) APIHandler(ctx context.Context) payV1.PaymentServiceServer {
	if d.apiHandler == nil {
		d.apiHandler = api.NewApi(d.PaymentService(ctx))
		logger.Info(ctx, "✅ Payment API handler initialized")
	}
	return d.apiHandler
}
