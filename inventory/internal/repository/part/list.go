package part

import (
	"github.com/Ilia9531/microservices-warehouse/inventory/internal/model"
	repoConverter "github.com/Ilia9531/microservices-warehouse/inventory/internal/repository/converter"
	repoModel "github.com/Ilia9531/microservices-warehouse/inventory/internal/repository/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) List(ctx context.Context, filter1 *model.PartsFilter) ([]*model.Part, error) {
	query := bson.M{}

	//fmt.Printf("Я List из репо, вот какой пришел filter: %v, коллекция: %v, датабаза: %v\n",
	//	filter1,
	//	r.collection.Name(),
	//	r.collection.Database().Name(),
	//)

	filter := repoConverter.FilterToRepo(filter1)

	// 1. Фильтр по UUID: $in = логическое ИЛИ внутри массива
	if len(filter.Uuids) > 0 {
		query["uuid"] = bson.M{"$in": filter.Uuids}
	}

	// 2. Фильтр по именам: $in
	if len(filter.Names) > 0 {
		query["name"] = bson.M{"$in": filter.Names}
	}

	// 3. Фильтр по категориям: конвертируем enum в строки, затем $in
	if len(filter.Category) > 0 {
		cats := make([]string, 0, len(filter.Category))
		for _, c := range filter.Category {
			cats = append(cats, string(c))
		}
		query["category"] = bson.M{"$in": cats}
	}

	// 4. Фильтр по странам производителя: вложенное поле через точку
	// В модели БД это: "manufacturer": { "country": "Russia" }
	if len(filter.ManufacturerCountries) > 0 {
		query["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}

	// 5. Фильтр по тегам: $in (любой из тегов совпадает)
	// Если нужно "все теги сразу", используй $all вместо $in
	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}

	// Выполняем запрос
	cursor, err := r.collection.Find(ctx, query)
	if err != nil {
		//fmt.Printf("Я list из репо, получена ошибка: %v, вот ее тип: %T\n", err, err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var repoModels []repoModel.Part
	if err = cursor.All(ctx, &repoModels); err != nil {
		return nil, err
	}

	listRes := make([]*model.Part, 0, len(repoModels))
	for _, item := range repoModels {
		x := repoConverter.PartToDomain(&item)
		listRes = append(listRes, x)
	}
	//fmt.Printf("Я list из репо, запрос list отправлен, получен listRes: %v\n", listRes)

	return listRes, nil
}
