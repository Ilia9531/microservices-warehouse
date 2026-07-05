package converter

import (
	"github.com/Ilia9531/microservices-warehouse/inventory/internal/model"
	inventoryV1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/inventory/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func PartToProto(p *model.Part) *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      CategoryToProto(p.Category),
		Dimensions:    DimensionsToProto(p.Dimensions),
		Manufacturer:  ManufacturerToProto(p.Manufacturer),
		Tags:          p.Tags,
		Metadata:      MetadataToProto(p.Metadata),
		CreatedAt:     timestamppb.New(p.CreatedAt),
		UpdatedAt:     timestamppb.New(p.UpdatedAt),
	}
}

func PartsFilterFromProto(f *inventoryV1.PartsFilter) *model.PartsFilter {
	if f == nil {
		return &model.PartsFilter{}
	}

	categories := make([]model.Category, 0, len(f.Categories))
	for _, c := range f.Categories {
		categories = append(categories, CategoryFromProto(c))
	}

	return &model.PartsFilter{
		Uuids:                 f.Uuids,
		Names:                 f.Names,
		Category:              categories,
		ManufacturerCountries: f.ManufacturerCountries,
		Tags:                  f.Tags,
	}
}

// PartsToProto — конвертирует список деталей из domain в proto
func PartsToProto(parts []*model.Part) []*inventoryV1.Part {
	if parts == nil {
		return nil
	}
	result := make([]*inventoryV1.Part, 0, len(parts))
	for _, p := range parts {
		result = append(result, PartToProto(p))
	}
	return result
}

func CategoryToProto(c model.Category) inventoryV1.Category {
	switch c {
	case model.CategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNKNOWN
	}
}

func CategoryFromProto(c inventoryV1.Category) model.Category {
	switch c {
	case inventoryV1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryV1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryV1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}

func DimensionsToProto(d model.Dimensions) *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func ManufacturerToProto(m model.Manufacturer) *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

// Для Metadata (map[string]MetaValue)
func MetadataToProto(m map[string]model.MetaValue) map[string]*inventoryV1.Value {
	if m == nil {
		return nil
	}
	result := make(map[string]*inventoryV1.Value, len(m))
	for k, v := range m {
		result[k] = MetaValueToProto(v)
	}
	return result
}

func MetaValueToProto(v model.MetaValue) *inventoryV1.Value {
	protoVal := &inventoryV1.Value{}
	if v.StringValue != nil {
		protoVal.Value = &inventoryV1.Value_StringValue{StringValue: *v.StringValue}
	} else if v.Int64Value != nil {
		protoVal.Value = &inventoryV1.Value_Int64Value{Int64Value: *v.Int64Value}
	} else if v.DoubleValue != nil {
		protoVal.Value = &inventoryV1.Value_DoubleValue{DoubleValue: *v.DoubleValue}
	} else if v.BoolValue != nil {
		protoVal.Value = &inventoryV1.Value_BoolValue{BoolValue: *v.BoolValue}
	}
	return protoVal
}
