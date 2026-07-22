package converter

//import (
//	"github.com/Ilia9531/microservices-warehouse/iam/internal/model"
//	repoModel "github.com/Ilia9531/microservices-warehouse/iam/internal/repository/model"
//)
//
//// ToRepoSession конвертирует доменную сессию в модель репозитория (Redis).
//func ToRepoSession(s *model.Session) *repoModel.Session {
//	if s == nil {
//		return nil
//	}
//	return &repoModel.Session{
//		UUID:     s.UUID,
//		UserUUID: s.UserUUID,
//	}
//}
//
//// FromRepoSession конвертирует модель репозитория в доменную сессию.
//func FromRepoSession(rs *repoModel.Session) *model.Session {
//	if rs == nil {
//		return nil
//	}
//	return &model.Session{
//		UUID:     rs.UUID,
//		UserUUID: rs.UserUUID,
//	}
//}
