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

func (a *Api) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := a.userService.GetUser(ctx, req.UserUuid)
	if err != nil {
		logger.Error(ctx, "failed to GetUser",
			zap.String("user_uuid", req.UserUuid),
			zap.Error(err),
		)
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
	}

	// internal/converter: Domain Model → Proto Message
	return &userv1.GetUserResponse{
		User: converter.ToProtoUser(user),
	}, nil
}
