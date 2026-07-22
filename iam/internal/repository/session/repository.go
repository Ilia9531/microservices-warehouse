package session

import (
	"fmt"
	"time"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/repository"
	"github.com/Ilia9531/microservices-warehouse/platform/pkg/cache"
)

var _ repository.SessionRepository = (*redisCache)(nil)

const cacheKeyPrefix = "iam:session:"

type redisCache struct {
	cache      cache.RedisClient
	sessionTTL time.Duration
}

func NewRepository(cache cache.RedisClient, sessionTTL time.Duration) *redisCache {
	return &redisCache{
		cache:      cache,
		sessionTTL: sessionTTL,
	}
}

func (r *redisCache) getCacheKey(sessionUUID string) string {
	return fmt.Sprintf("%s%s", cacheKeyPrefix, sessionUUID)
}

func (r *redisCache) getUserSessionsKey(userUUID string) string {
	return fmt.Sprintf("iam:user:sessions:%s", userUUID)
}
