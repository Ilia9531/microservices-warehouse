package v1

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Ilia9531/microservices-warehouse/inventory/internal/converter"
	"github.com/Ilia9531/microservices-warehouse/inventory/internal/model"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/logger"
	inventoryV1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/inventory/v1"
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
