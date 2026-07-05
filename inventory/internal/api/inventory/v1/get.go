package v1

import (
	"Jopa/inventory/internal/converter"
	"Jopa/inventory/internal/model"
	"Jopa/platform/pkg/logger"
	inventoryV1 "Jopa/shared/pkg/proto/inventory/v1"
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {

	part, err := a.InventoryService.GetPart(ctx, req.GetUuid())
	if err != nil {
		logger.Error(ctx, "failed to get part",
			zap.String("uuid", req.Uuid),
			zap.Error(err),
		)
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part:%s not found", req.GetUuid())
		}
		return nil, err
	}

	return &inventoryV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
