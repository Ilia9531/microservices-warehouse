package part

import (
	"github.com/Ilia9531/microservices-warehouse/inventory/internal/model"
	repoConverter "github.com/Ilia9531/microservices-warehouse/inventory/internal/repository/converter"
	RepoModel "github.com/Ilia9531/microservices-warehouse/inventory/internal/repository/model"
	"context"
	"errors"
	"fmt"

	gUuid "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) Get(ctx context.Context, uuid string) (*model.Part, error) {
	_, err := gUuid.Parse(uuid)
	if err != nil {
		return nil, model.ErrValidation
	}
	//fmt.Printf("Я get из репо, вот какой пришел uuid: %v, коллекция: %v, датабаза: %v\n",
	//	uuid,
	//	r.collection.Name(),
	//	r.collection.Database().Name(),
	//)
	var dbPart RepoModel.Part
	// Фильтр: ищем документ, где поле uuid равно переданному значению
	filter := bson.M{"uuid": uuid}

	// FindOne возвращает курсор на один документ. Decode распаковывает BSON в Go-структуру.
	err = r.collection.FindOne(ctx, filter).Decode(&dbPart)

	if err != nil {
		// MongoDB возвращает mongo.ErrNoDocuments, если документ не найден
		fmt.Printf("Я get из репо, получена ошибка: %v, вот ее тип: %T\n", err, err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrPartNotFound
		}
		return nil, err
	}
	fmt.Printf("Я get из репо, запрос get отправлен, получен dbPart: %v\n", dbPart)
	return repoConverter.PartToDomain(&dbPart), nil
}
