package v1

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/converter"
	"github.com/Ilia9531/microservices-warehouse/iam/internal/model"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/logger"
	userv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/user/v1"
)

func (a *Api) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	// Сервис возвращает UUID созданного пользователя
	modelReg := converter.ToDomainRegisterRequest(req)
	userUUID, err := a.userService.Register(ctx, modelReg)
	if err != nil {
		logger.Error(ctx, "failed to Register",
			zap.String("session_uuid", req.Login),
			zap.Error(err),
		)
		if errors.Is(err, model.ErrUserAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "User already exists")
		}
	}

	return &userv1.RegisterResponse{UserUuid: userUUID}, nil
}
