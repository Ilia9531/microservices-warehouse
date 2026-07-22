package converter

import (
	"encoding/json"
	"fmt"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/model"
	repoModel "github.com/Ilia9531/microservices-warehouse/iam/internal/repository/model"
	commonv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/common/v1"
)

// ToRepoUser конвертирует доменную модель в модель репозитория (PG).
func ToRepoUser(u *model.User) *repoModel.User {
	if u == nil {
		return nil
	}
	// Сериализуем notification_methods в JSON
	methodsJSON, _ := json.Marshal(u.NotificationMethods)
	return &repoModel.User{
		UUID:                u.UUID,
		Login:               u.Login,
		PasswordHash:        u.PasswordHash,
		Email:               u.Email,
		NotificationMethods: methodsJSON,
	}
}

// FromRepoUser конвертирует модель репозитория в доменную модель.
func FromRepoUser(ru *repoModel.User) (*model.User, error) {
	if ru == nil {
		return nil, fmt.Errorf("user is nil")
	}
	var methods []*commonv1.NotificationMethod
	if len(ru.NotificationMethods) > 0 {
		err := json.Unmarshal(ru.NotificationMethods, &methods)
		if err != nil {
			return nil, err
		}
	}

	return &model.User{
		UUID:                ru.UUID,
		Login:               ru.Login,
		PasswordHash:        ru.PasswordHash,
		Email:               ru.Email,
		NotificationMethods: methods,
	}, nil
}
