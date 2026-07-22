package session

import (
	"context"
)

// AddSessionToUserSet индексирует сессию, добавляя её UUID во множество сессий пользователя.
// Это позволяет в будущем реализовать функции вроде "выйти со всех устройств" или
// проверки количества одновременных сессий.
func (r *redisCache) AddSessionToUserSet(ctx context.Context, userUUID, sessionUUID string) error {
	key := r.getUserSessionsKey(userUUID)
	return r.cache.SAdd(ctx, key, sessionUUID)
}
