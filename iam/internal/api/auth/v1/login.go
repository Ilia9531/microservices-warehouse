package v1

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/model"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/logger"
	authv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/auth/v1"
)

func (a *Api) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	// Сервис проверяет пароль, создаёт сессию в Redis и возвращает SessionUUID
	sessionUUID, err := a.authService.Login(ctx, req.Login, req.Password)
	if err != nil {
		logger.Error(ctx, "failed to Login",
			zap.String("login", req.Login),
			zap.Error(err),
		)
		if errors.Is(err, model.ErrInvalidLogin) {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid login or password")
		}

	}

	return &authv1.LoginResponse{SessionUuid: sessionUUID}, nil
}
