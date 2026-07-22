package session

import (
	"context"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/model"
	"github.com/Ilia9531/microservices-warehouse/iam/internal/repository/converter"
	repoModel "github.com/Ilia9531/microservices-warehouse/iam/internal/repository/model"
)

func (r *redisCache) GetByUUID(ctx context.Context, uuid string) (*model.Session, error) {
	if uuid == "" {
		return nil, model.ErrSessionInvalid
	}
	cacheKey := r.getCacheKey(uuid)

	// Получаем Hash из Redis
	values, err := r.cache.HGetAll(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return nil, model.ErrSessionNotFound
		}
		return nil, err
	}

	if len(values) == 0 {
		return nil, model.ErrSessionNotFound
	}

	// Десериализуем в Redis-view
	var redisView repoModel.Session
	if err = redigo.ScanStruct(values, &redisView); err != nil {
		return nil, fmt.Errorf("fail on redigo.ScanStruct: %w", err)
	}

	// Конвертируем в доменную модель репозитория
	return converter.FromRepoSession(&redisView), nil
}
