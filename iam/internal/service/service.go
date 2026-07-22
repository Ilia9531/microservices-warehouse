package service

import (
	"context"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/model"
)

type UserService interface {
	GetUser(ctx context.Context, uuid string) (*model.User, error)
	Register(ctx context.Context, req *model.RegisterRequest) (string, error)
}

type AuthService interface {
	Login(ctx context.Context, login, password string) (sessionUuid string, err error)
	Whoami(ctx context.Context, sessionUuid string) (*model.Whoami, error)
}
