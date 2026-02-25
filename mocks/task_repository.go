package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tdatIT/backend-go/internal/domain/models"
	taskrepo "github.com/tdatIT/backend-go/internal/infras/repository/task"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, item *models.Task) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockTaskRepository) FindByID(ctx context.Context, id uint64) (*models.Task, error) {
	args := m.Called(ctx, id)
	var result *models.Task
	if args.Get(0) != nil {
		result = args.Get(0).(*models.Task)
	}
	return result, args.Error(1)
}

func (m *MockTaskRepository) FindAllBy(ctx context.Context, params *taskrepo.GetListParams) ([]*models.Task, int64, error) {
	args := m.Called(ctx, params)
	var results []*models.Task
	if args.Get(0) != nil {
		results = args.Get(0).([]*models.Task)
	}
	var total int64
	if args.Get(1) != nil {
		total = args.Get(1).(int64)
	}
	return results, total, args.Error(2)
}

func (m *MockTaskRepository) Update(ctx context.Context, item *models.Task) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
