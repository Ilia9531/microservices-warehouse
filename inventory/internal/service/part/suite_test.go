package part_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Ilia9531/microservices-warehouse/inventory/internal/model"
	"github.com/Ilia9531/microservices-warehouse/inventory/internal/repository/mocks"
	"github.com/Ilia9531/microservices-warehouse/inventory/internal/service/part"
)

type PartServiceTestSuite struct {
	suite.Suite
	ctx     context.Context
	repo    *mocks.InventoryRepository
	service *part.Service
}

func (s *PartServiceTestSuite) SetupTest() {
	// Создаём мок репозитория
	s.repo = mocks.NewInventoryRepository(s.T())

	// Создаём сервис с моком
	s.service = part.NewService(s.repo)
}

// TestPartServiceSuite — точка входа для запуска всех тестов набора.
func TestPartServiceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(PartServiceTestSuite))
}

// Хелпер: создание тестовой детали.
func (s *PartServiceTestSuite) createTestPart() *model.Part {
	return &model.Part{
		Uuid:     "test-part-uuid",
		Name:     "Test Part",
		Price:    99.99,
		Category: model.CategoryEngine,
	}
}
