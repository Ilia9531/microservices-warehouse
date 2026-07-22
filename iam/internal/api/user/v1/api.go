package v1

import (
	"github.com/Ilia9531/microservices-warehouse/iam/internal/service"
	userv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/user/v1"
)

type Api struct {
	userv1.UnimplementedUserServiceServer
	userService service.UserService
}

func NewAPI(userService service.UserService) *Api {
	return &Api{userService: userService}
}
