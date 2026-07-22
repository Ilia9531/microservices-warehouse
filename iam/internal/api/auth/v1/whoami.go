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

func (a *Api) Whoami(ctx context.Context, req *authv1.WhoamiRequest) (*authv1.WhoamiResponse, error) {
	// session_uuid передаётся в теле запроса (grpcurl test)
	user, err := a.authService.Whoami(ctx, req.SessionUuid)
	if err != nil {
		logger.Error(ctx, "failed to Whoami",
			zap.String("session_uuid", req.SessionUuid),
			zap.Error(err),
		)
		if errors.Is(err, model.ErrSessionNotFound) {
			return nil, status.Errorf(codes.NotFound, "Session not found")
		}
		if errors.Is(err, model.ErrSessionExpired) {
			return nil, status.Errorf(codes.Unauthenticated, "Session was expired")
		}
		if errors.Is(err, model.ErrSessionInvalid) {
			return nil, status.Errorf(codes.Unauthenticated, "Session is invalid")
		}
	}

	return &authv1.WhoamiResponse{
		UserUuid: user.UUID,
		Login:    user.Login,
		Email:    user.Email,
	}, nil
}
