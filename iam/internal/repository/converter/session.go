package converter

import (
	"github.com/Ilia9531/microservices-warehouse/iam/internal/model"
	repoModel "github.com/Ilia9531/microservices-warehouse/iam/internal/repository/model"
)

func ToRepoSession(s *model.Session) *repoModel.Session {
	if s == nil {
		return nil
	}
	return &repoModel.Session{
		UUID:     s.UUID,
		UserUUID: s.UserUUID,
	}
}

func FromRepoSession(s *repoModel.Session) *model.Session {
	if s == nil {
		return nil
	}
	return &model.Session{
		UUID:     s.UUID,
		UserUUID: s.UserUUID,
	}
}
