package part_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Jopa/inventory/internal/model"
)

// TestListParts_Success тест на успешное получение списка деталей.
func (s *PartServiceTestSuite) TestListParts_Success() {
	ctx := s.ctx

	// Создаём тестовые данные
	parts := []*model.Part{
		s.createTestPart(),
		{
			Uuid:     "part-2",
			Name:     "Part 2",
			Price:    199.99,
			Category: model.CategoryFuel,
		},
	}
	// Мокируем успешное получение списка из репозитория
	s.repo.EXPECT().
		List(ctx, mock.Anything). // filter может быть nil или любым
		Return(parts, nil).
		Once()
	// Вызываем сервис (без фильтра — все детали)
	result, err := s.service.ListParts(ctx, nil)

	// Проверяем результат
	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Len(result, 2)
	s.Equal("test-part-uuid", result[0].Uuid)
	s.Equal("part-2", result[1].Uuid)
}

// TestListParts_WithFilter тест на получение списка с фильтром.
func (s *PartServiceTestSuite) TestListParts_WithFilter() {
	ctx := s.ctx

	// Фильтр по UUID
	filter := &model.PartsFilter{
		Names: []string{"Test Part"},
	}

	// Создаём тестовые данные (только 1 деталь соответствует фильтру)
	parts := []*model.Part{
		s.createTestPart(),
	}

	// Мокируем получение списка с фильтром
	s.repo.EXPECT().
		List(ctx, filter).
		Return(parts, nil).
		Once()

	// Вызываем сервис с фильтром
	result, err := s.service.ListParts(ctx, filter)

	// Проверяем результат
	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Len(result, 1)
	s.Equal("Test Part", result[0].Name)
}

// TestListParts_EmptyResult тест на пустой результат.
func (s *PartServiceTestSuite) TestListParts_EmptyResult() {
	ctx := s.ctx

	// Мокируем пустой список (не ошибка!)
	s.repo.EXPECT().
		List(ctx, mock.Anything).
		Return([]*model.Part{}, nil).
		Once()

	// Вызываем сервис
	result, err := s.service.ListParts(ctx, nil)

	// Проверяем результат
	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Len(result, 0) // пустой срез — это нормально
}

// TestListParts_RepositoryError тест на ошибку репозитория.
func (s *PartServiceTestSuite) TestListParts_RepositoryError() {
	ctx := s.ctx

	// Мокируем ошибку репозитория
	s.repo.EXPECT().
		List(ctx, mock.Anything).
		Return(nil, assert.AnError).
		Once()

	// Вызываем сервис
	result, err := s.service.ListParts(ctx, nil)

	// Проверяем результат
	s.Require().Error(err)
	s.Require().Nil(result)
}
