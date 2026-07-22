package session

import (
	"context"
	"time"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/model"
	"github.com/Ilia9531/microservices-warehouse/iam/internal/repository/converter"
)

func (r *redisCache) Create(ctx context.Context, session *model.Session, _ time.Duration) error {
	// TTL фиксированный, игнорируем переданный параметр (согласно требованию)
	// Формируем ключ в Redis
	cacheKey := r.getCacheKey(session.UUID)

	// Конвертируем в Redis-view с тегами redis:"..."
	redisView := converter.ToRepoSession(session)

	// Сохраняем как Hash
	if err := r.cache.HashSet(ctx, cacheKey, redisView); err != nil {
		return err
	}

	// Устанавливаем TTL
	if err := r.cache.Expire(ctx, cacheKey, r.sessionTTL); err != nil {
		return err
	}

	// Индексируем сессию по пользователю (для будущих фич)

	return r.AddSessionToUserSet(ctx, session.UserUUID, session.UUID)
}
