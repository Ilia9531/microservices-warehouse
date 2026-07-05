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

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {

	//if req.GetFilter() != nil {
	//	fmt.Printf("Я ListParts из api, вот параметры из фильтра: names=%v, categories=%v, tags=%v\n",
	//		req.GetFilter().GetNames(),
	//		req.GetFilter().GetCategories(),
	//		req.GetFilter().GetTags())
	//}

	filter := converter.PartsFilterFromProto(req.GetFilter())

	parts, err := a.InventoryService.ListParts(ctx, filter)
	if err != nil {
		logger.Error(ctx, "failed to get parts",
			zap.Error(err),
		)
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "Filter:%s not found", req.GetFilter())
		}
		return nil, err
	}
	//if len(parts) > 0 {
	//	fmt.Printf("Я listParts из api, обратно отправлены: uuid=%v, name=%v, price=%v\n",
	//		parts[0].Uuid, parts[0].Name, parts[0].Price,
	//	)
	//} else {
	//	fmt.Printf("Я listParts из api, вернулся пустым: %v\n", parts)
	//}

	return &inventoryV1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}, nil
}
