//go:build integration

package integration

import (
	repoModel "Jopa/inventory/internal/repository/model"
	"Jopa/platform/pkg/testcontainers/app"
	"Jopa/platform/pkg/testcontainers/mongo"
	"Jopa/platform/pkg/testcontainers/network"
	invV1 "Jopa/shared/pkg/proto/inventory/v1"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestEnvironment хранит ссылки на контейнеры тестового окружения
type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}

// InsertTestPart — вставляет тестовую деталь в коллекцию MongoDB и возвращает её UUID
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := gofakeit.UUID()
	now := time.Now()

	partDoc := bson.M{
		"uuid":           partUUID,
		"name":           gofakeit.ProductName(),
		"description":    gofakeit.Paragraph(3, 3, 10, " "),
		"price":          gofakeit.Price(10, 10000),
		"stock_quantity": gofakeit.Number(1, 1000),
		"category":       gofakeit.RandomString([]string{"ENGINE", "FUEL", "PORTHOLE", "WING"}),
		"dimensions": bson.M{
			"length": gofakeit.Float64Range(10, 500),
			"width":  gofakeit.Float64Range(10, 500),
			"height": gofakeit.Float64Range(10, 500),
			"weight": gofakeit.Float64Range(1, 1000),
		},
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
			"website": gofakeit.URL(),
		},
		"tags": []string{gofakeit.Word(), gofakeit.Word()},
		"metadata": bson.M{
			"serial": bson.M{
				"string_value": gofakeit.Word(), // ← правильная структура
			},
			"batch": bson.M{
				"int64_value": int64(gofakeit.Number(1, 100)), // ← правильная структура
			},
		},
		"created_at": primitive.NewDateTimeFromTime(now),
		"updated_at": primitive.NewDateTimeFromTime(now),
	}

	databaseName := getEnvWithFallback(mongoDatabaseKey, "inventory")
	collection := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName)

	//лог

	res, err := collection.InsertOne(ctx, partDoc)
	if err != nil {
		return "", fmt.Errorf("failed to insert test part: %w", err)
	}
	//лог2
	fmt.Printf("✅ InsertTestPart: inserted ID=%v\n", res.InsertedID)
	//лог3

	return partUUID, nil
}

// InsertTestPartWithData — вставляет тестовую деталь с заданными данными
func (env *TestEnvironment) InsertTestPartWithData(ctx context.Context, part *invV1.Part) (string, error) {
	partUUID := part.GetUuid()
	if partUUID == "" {
		partUUID = gofakeit.UUID()
	}
	now := time.Now()

	partDoc := bson.M{
		"uuid":           partUUID,
		"name":           part.GetName(),
		"description":    part.GetDescription(),
		"price":          part.GetPrice(),
		"stock_quantity": part.GetStockQuantity(),
		"category":       part.GetCategory().String(),
		"dimensions": bson.M{
			"length": part.GetDimensions().GetLength(),
			"width":  part.GetDimensions().GetWidth(),
			"height": part.GetDimensions().GetHeight(),
			"weight": part.GetDimensions().GetWeight(),
		},
		"manufacturer": bson.M{
			"name":    part.GetManufacturer().GetName(),
			"country": part.GetManufacturer().GetCountry(),
			"website": part.GetManufacturer().GetWebsite(),
		},
		"tags":       part.GetTags(),
		"metadata":   convertMetadataToBSON(part.GetMetadata()),
		"created_at": primitive.NewDateTimeFromTime(now),
		"updated_at": primitive.NewDateTimeFromTime(now),
	}

	databaseName := getEnvWithFallback(mongoDatabaseKey, "inventory")
	collection := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName)
	//лог

	_, err := collection.InsertOne(ctx, partDoc)
	if err != nil {
		return "", fmt.Errorf("failed to insert test part: %w", err)
	}

	return partUUID, nil
}

// GetTestPart — возвращает тестовую деталь для создания
func (env *TestEnvironment) GetTestPart() *invV1.Part {
	now := time.Now()
	return &invV1.Part{
		Uuid:          gofakeit.UUID(),
		Name:          "Test Booster MK-III",
		Description:   "High-performance booster for interplanetary travel",
		Price:         15000.50,
		StockQuantity: 42,
		Category:      invV1.Category_CATEGORY_ENGINE,
		Dimensions: &invV1.Dimensions{
			Length: 250.5,
			Width:  80.0,
			Height: 80.0,
			Weight: 1200.0,
		},
		Manufacturer: &invV1.Manufacturer{
			Name:    "SpaceTech Industries",
			Country: "Germany",
			Website: "https://spacetech.example",
		},
		Tags: []string{"booster", "engine", "high-performance"},
		Metadata: map[string]*invV1.Value{
			"serial": {Value: &invV1.Value_StringValue{StringValue: "SN-2026-001"}},
			"batch":  {Value: &invV1.Value_Int64Value{Int64Value: 42}},
		},
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}
}

// GetTestPartFilter — возвращает тестовый фильтр для ListParts
func (env *TestEnvironment) GetTestPartFilter() *invV1.PartsFilter {
	return &invV1.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{"Test Booster"},
		Categories:            []invV1.Category{invV1.Category_CATEGORY_ENGINE},
		ManufacturerCountries: []string{"Germany"},
		Tags:                  []string{"booster"},
	}
}

// ClearPartsCollection — удаляет все записи из коллекции parts
func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	databaseName := getEnvWithFallback(mongoDatabaseKey, "inventory")

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to clear parts collection: %w", err)
	}
	return nil
}

// getEnvWithFallback возвращает значение переменной окружения или значение по умолчанию
func getEnvWithFallback(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func convertMetadataToBSON(metadata map[string]*invV1.Value) map[string]repoModel.MetaValue {
	result := make(map[string]repoModel.MetaValue)
	for k, v := range metadata {
		if v == nil {
			continue
		}

		metaValue := repoModel.MetaValue{}
		switch val := v.Value.(type) {
		case *invV1.Value_StringValue:
			metaValue.StringValue = &val.StringValue
		case *invV1.Value_Int64Value:
			metaValue.Int64Value = &val.Int64Value
		case *invV1.Value_DoubleValue:
			metaValue.DoubleValue = &val.DoubleValue
		case *invV1.Value_BoolValue:
			metaValue.BoolValue = &val.BoolValue
		}
		result[k] = metaValue
	}
	return result
}
