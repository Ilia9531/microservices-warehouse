package converter

import (
	"github.com/Ilia9531/microservices-warehouse/inventory/internal/model"
	repoModel "github.com/Ilia9531/microservices-warehouse/inventory/internal/repository/model"
)

func PartToDomain(p *repoModel.Part) *model.Part {
	return &model.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      model.Category(p.Category),
		Dimensions:    model.Dimensions(p.Dimensions),
		Manufacturer:  model.Manufacturer(p.Manufacturer),
		Tags:          p.Tags,
		Metadata:      MetadataToDomain(p.Metadata),
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func PartToRepo(p *model.Part) *repoModel.Part {
	return &repoModel.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      string(p.Category),
		Dimensions:    repoModel.Dimensions(p.Dimensions),
		Manufacturer:  repoModel.Manufacturer(p.Manufacturer),
		Tags:          p.Tags,
		Metadata:      DomainToMetadata(p.Metadata),
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

// CategoryToDomain конвертирует repo Category в domain Category
func CategoryToDomain(c repoModel.Category) model.Category {
	return model.Category(c)
}

// FilterToDomain конвертирует repo-фильтр в domain-фильтр
func FilterToDomain(p *repoModel.PartsFilter) *model.PartsFilter {
	// Конвертируем каждый элемент слайса Category
	categories := make([]model.Category, 0, len(p.Category))
	for _, c := range p.Category {
		categories = append(categories, CategoryToDomain(c))
	}

	return &model.PartsFilter{
		Uuids:                 p.Uuids,
		Names:                 p.Names,
		Category:              categories, // ← используем конвертированный слайс
		ManufacturerCountries: p.ManufacturerCountries,
		Tags:                  p.Tags,
	}
}

// CategoryToRepo конвертирует domain Category в repo Category
func CategoryToRepo(c model.Category) repoModel.Category {
	return repoModel.Category(c)
}

// FilterToRepo конвертирует domain-фильтр в repo-фильтр
func FilterToRepo(p *model.PartsFilter) *repoModel.PartsFilter {
	categories := make([]repoModel.Category, 0, len(p.Category))
	for _, c := range p.Category {
		categories = append(categories, CategoryToRepo(c))
	}

	return &repoModel.PartsFilter{
		Uuids:                 p.Uuids,
		Names:                 p.Names,
		Category:              categories,
		ManufacturerCountries: p.ManufacturerCountries,
		Tags:                  p.Tags,
	}
}

// -----------------------------------------------------------------------------
// Нейро-заглушка
func MetaValueToDomain(v repoModel.MetaValue) model.MetaValue {
	return model.MetaValue{
		StringValue: v.StringValue,
		Int64Value:  v.Int64Value,
		DoubleValue: v.DoubleValue,
		BoolValue:   v.BoolValue,
	}
}

// Нейро-заглушка
func MetadataToDomain(repoMeta map[string]repoModel.MetaValue) map[string]model.MetaValue {
	domainMeta := make(map[string]model.MetaValue, len(repoMeta))
	for k, v := range repoMeta {
		domainMeta[k] = MetaValueToDomain(v)
	}
	return domainMeta
}

// Нейро-заглушка
func DomainToMetaValue(v model.MetaValue) repoModel.MetaValue {
	return repoModel.MetaValue{
		StringValue: v.StringValue,
		Int64Value:  v.Int64Value,
		DoubleValue: v.DoubleValue,
		BoolValue:   v.BoolValue,
	}
}

// Нейро-заглушка
func DomainToMetadata(Meta map[string]model.MetaValue) map[string]repoModel.MetaValue {
	domainMeta := make(map[string]repoModel.MetaValue, len(Meta))
	for k, v := range Meta {
		domainMeta[k] = DomainToMetaValue(v)
	}
	return domainMeta
}
