package user

import (
	"github.com/Ilia9531/microservices-warehouse/iam/internal/repository"
	def "github.com/Ilia9531/microservices-warehouse/iam/internal/service"
)

var _ def.UserService = (*Service)(nil)

type Service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) *Service {
	return &Service{userRepo: userRepo}
}
