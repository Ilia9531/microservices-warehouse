package repository

import (
	"context"
	"time"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/model"
)

// UserRepository управляет данными пользователей в PostgreSQL.
// Реализация находится в internal/repository/user/repository.go
type UserRepository interface {
	// Create сохраняет нового пользователя в БД.
	Create(ctx context.Context, user *model.User) error

	// GetByUUID ищет пользователя по UUID.
	// Возвращает ошибку, если запись не найдена.
	GetByUUID(ctx context.Context, uuid string) (*model.User, error)

	// GetByLogin ищет пользователя по логину.
	// Возвращает ошибку, если запись не найдена.
	GetByLogin(ctx context.Context, login string) (*model.User, error)
}

// SessionRepository управляет активными сессиями в Redis.
// Реализация находится в internal/repository/session/repository.go
type SessionRepository interface {
	// Create сохраняет сессию в Redis с указанным временем жизни (TTL).
	Create(ctx context.Context, session *model.Session, ttl time.Duration) error

	// GetByUUID получает сессию по идентификатору.
	// Возвращает ошибку, если сессия не найдена или срок её действия истёк.
	GetByUUID(ctx context.Context, uuid string) (*model.Session, error)
}
