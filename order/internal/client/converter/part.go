package converter

import (
	"github.com/Ilia9531/microservices-warehouse/order/internal/model"
	inventoryv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/inventory/v1"
)

func ToDomainPart(protoPart *inventoryv1.Part) *model.Part{
	if protoPart == nil {
		return nil
	}
	return &model.Part{
		UUID: protoPart.Uuid,
		Name: protoPart.Name,
		Price: protoPart.Price,
	}
}

func ToDomainParts(protoParts []*inventoryv1.Part) []*model.Part {
	if protoParts == nil {
		return nil
	}
	domainParts := make([]*model.Part, 0, len(protoParts))
	for _, protoPart := range protoParts {
		domainParts = append(domainParts, ToDomainPart(protoPart))
	}
	return domainParts
}
