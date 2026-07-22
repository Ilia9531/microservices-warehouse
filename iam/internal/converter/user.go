package converter

import (
	"github.com/Ilia9531/microservices-warehouse/iam/internal/model"
	"github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/common/v1"
	userv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/user/v1"
)

// ToProtoUser конвертирует доменную модель в proto-сообщение.
func ToProtoUser(u *model.User) *commonv1.User {
	if u == nil {
		return nil
	}
	return &commonv1.User{
		Uuid:                u.UUID,
		Login:               u.Login,
		Email:               u.Email,
		NotificationMethods: u.NotificationMethods,
	}
}

func FromProtoUser(u *commonv1.User) *model.User {
	if u == nil {
		return nil
	}
	return &model.User{
		UUID:                u.Uuid,
		Login:               u.Login,
		Email:               u.Email,
		NotificationMethods: u.NotificationMethods,
	}
}

// ToDomainRegisterRequest преобразует proto-запрос регистрации в доменную модель.
func ToDomainRegisterRequest(req *userv1.RegisterRequest) *model.RegisterRequest {
	if req == nil {
		return nil
	}
	return &model.RegisterRequest{
		Login:               req.Login,
		Password:            req.Password,
		Email:               req.Email,
		NotificationMethods: req.NotificationMethods,
	}
}
