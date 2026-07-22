package v1

import (
	"github.com/Ilia9531/microservices-warehouse/iam/internal/service"
	authv1 "github.com/Ilia9531/microservices-warehouse/shared/pkg/proto/auth/v1"
)

type Api struct {
	authv1.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewAPI(authService service.AuthService) *Api {
	return &Api{authService: authService}
}
