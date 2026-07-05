package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	invV1API "github.com/Ilia9531/microservices-warehouse/inventory/internal/api/inventory/v1"
	"github.com/Ilia9531/microservices-warehouse/inventory/internal/config"
	"github.com/Ilia9531/microservices-warehouse/inventory/internal/repository"
	invRep "github.com/Ilia9531/microservices-warehouse/inventory/internal/repository/part"
	"github.com/Ilia9531/microservices-warehouse/inventory/internal/service"
	PartServ "github.com/Ilia9531/microservices-warehouse/inventory/internal/service/part"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/closer"
	invV1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/inventory/v1"
)

const partsCollectionName = "parts"

type diContainer struct {
	invApiV1      invV1.InventoryServiceServer
	invService    service.InventoryService
	invRepository repository.InventoryRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InvApiV1(ctx context.Context) invV1.InventoryServiceServer {
	if d.invApiV1 == nil {
		d.invApiV1 = invV1API.NewApi(d.PartService(ctx))
	}
	return d.invApiV1
}

func (d *diContainer) PartRepository(ctx context.Context) repository.InventoryRepository {
	if d.invRepository == nil {
		d.invRepository = invRep.NewRepository(d.MongoDBHandle(ctx), partsCollectionName)
	}
	return d.invRepository
}

func (d *diContainer) PartService(ctx context.Context) service.InventoryService {
	if d.invService == nil {
		d.invService = PartServ.NewService(d.PartRepository(ctx))
	}
	return d.invService
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}
	return d.mongoDBHandle
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}
		fmt.Printf("Пинг выполнен, ВОТ ПУТЬ МОНГИ: %s, и Database: %s\n",
			config.AppConfig().Mongo.URI(),
			config.AppConfig().Mongo.DatabaseName())
		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}
