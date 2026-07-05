package part_test

import (
	"github.com/stretchr/testify/assert"
	
	"Jopa/inventory/internal/model"
)

func (s *PartServiceTestSuite) TestGetPart_Success() {
	ctx := s.ctx
	expectPart := s.createTestPart()

	// Мокируем успешное получение из репозитория
	s.repo.EXPECT().
		Get(ctx, expectPart.Uuid).
		Return(expectPart, nil).
		Once()

	part, err := s.service.GetPart(ctx, expectPart.Uuid)

	s.Require().NoError(err)
	s.Require().NotNil(part)
	s.Equal(expectPart.Uuid, part.Uuid)
	s.Equal(expectPart.Name, part.Name)
	s.Equal(expectPart.Price, part.Price)
}

// TestGetPart_NotFound тест на отсутствие детали.
func (s *PartServiceTestSuite) TestGetPart_NotFound() {
	ctx := s.ctx

	// Мокируем ошибку NotFound
	s.repo.EXPECT().
		Get(ctx, "non-existent").
		Return(nil, model.ErrPartNotFound).
		Once()

	// Вызываем сервис
	part, err := s.service.GetPart(ctx, "non-existent")

	// Проверяем результат
	s.Require().Error(err)
	s.Require().Nil(part)
	s.ErrorIs(err, model.ErrPartNotFound)
}

// TestGetPart_EmptyUUID тест на пустой UUID.
func (s *PartServiceTestSuite) TestGetPart_EmptyUUID() {
	ctx := s.ctx

	// Вызываем сервис с пустым UUID
	part, err := s.service.GetPart(ctx, "")

	// Проверяем результат
	s.Require().Error(err)
	s.Require().Nil(part)
	s.ErrorIs(err, model.ErrValidation)
}

// TestGetPart_RepositoryError тест на ошибку репозитория.
func (s *PartServiceTestSuite) TestGetPart_RepositoryError() {
	ctx := s.ctx

	// Мокируем произвольную ошибку репозитория
	s.repo.EXPECT().
		Get(ctx, "test-uuid").
		Return(nil, assert.AnError).
		Once()

	// Вызываем сервис
	part, err := s.service.GetPart(ctx, "test-uuid")

	// Проверяем результат
	s.Require().Error(err)
	s.Require().Nil(part)
}
