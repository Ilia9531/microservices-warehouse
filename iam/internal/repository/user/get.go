package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/model"
	"github.com/Ilia9531/microservices-warehouse/iam/internal/repository/converter"
	repoModel "github.com/Ilia9531/microservices-warehouse/iam/internal/repository/model"
)

// GetByUUID ищет пользователя по UUID.
func (r *Repository) GetByUUID(ctx context.Context, uuid string) (*model.User, error) {
	if uuid == "" {
		return nil, model.ErrInvalidUuid
	}

	var repoUser repoModel.User
	err := r.db.QueryRow(ctx, `
		SELECT uuid, login, password_hash, email, notification_methods, created_at
		FROM users
		WHERE uuid = $1
	`, uuid).Scan(&repoUser.UUID, &repoUser.Login, &repoUser.PasswordHash, &repoUser.Email,
		&repoUser.NotificationMethods, &repoUser.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}
	user, err := converter.FromRepoUser(&repoUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByLogin ищет пользователя по логину.
func (r *Repository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	if login == "" {
		return nil, model.ErrInvalidLogin
	}

	var repoUser repoModel.User
	err := r.db.QueryRow(ctx, `
		SELECT uuid, login, password_hash, email, notification_methods, created_at
		FROM users
		WHERE login = $1
	`, login).Scan(&repoUser.UUID, &repoUser.Login, &repoUser.PasswordHash, &repoUser.Email,
		&repoUser.NotificationMethods, &repoUser.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	user, err := converter.FromRepoUser(&repoUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}
