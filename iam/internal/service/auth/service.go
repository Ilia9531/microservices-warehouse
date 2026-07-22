package auth

import (
	"time"

	"github.com/Ilia9531/microservices-warehouse/iam/internal/repository"
	def "github.com/Ilia9531/microservices-warehouse/iam/internal/service"
)

var _ def.AuthService = (*Service)(nil)

type Service struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	sessionTTL  time.Duration
}

func NewService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	sessionTTL time.Duration,
) *Service {
	return &Service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		sessionTTL:  sessionTTL,
	}
}
