package part

import (
	def "github.com/Ilia9531/microservices-warehouse/inventory/internal/repository"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database, name string) *repository {
	collection := db.Collection(name)

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(err)
	}

	repo := &repository{collection: collection}
	//закоменить при e2e тестах:
	//if err = repo.SeedIfEmpty(ctx); err != nil {
	//	log.Printf("⚠️ failed to seed data: %v", err)
	//}
	return repo
}

// seedIfEmpty вставляет тестовые данные, только если коллекция пуста
func (r *repository) SeedIfEmpty(ctx context.Context) error {
	// Проверяем, есть ли уже документы
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("⚠️ failed to count documents: %v\n", err)
		return err
	}
	if count > 0 {
		return err // Данные уже есть, не перезаписываем
	}

	// Вставляем сид-данные
	seedParts := initParts()
	if len(seedParts) == 0 {
		return err
	}

	_, err = r.collection.InsertMany(ctx, seedParts)
	if err != nil {
		log.Printf("⚠️ failed to seed initial data: %v", err)
		return err
	}

	log.Println("✅ Seeded initial inventory data")
	return nil
}
