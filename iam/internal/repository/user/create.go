package user

import (
	"context"
	"fmt"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/model"
	"github.com/Ilia9531/microservices-warehouse/iam/internal/repository/converter"
)

// Create сохраняет нового пользователя в PostgreSQL.
func (r *Repository) Create(ctx context.Context, user *model.User) error {
	err := user.Validate()
	if err != nil {
		return fmt.Errorf("validate RegisterRequest failed in repo layer: %w", err)
	}
	repoUser := converter.ToRepoUser(user)
	_, err = r.db.Exec(ctx, `
		INSERT INTO users (uuid, login, password_hash, email, notification_methods, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`, repoUser.UUID, repoUser.Login, repoUser.PasswordHash, repoUser.Email, repoUser.NotificationMethods)

	return err
}
