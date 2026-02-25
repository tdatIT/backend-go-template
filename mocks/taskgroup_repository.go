package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tdatIT/backend-go/internal/domain/models"
	taskgrouprepo "github.com/tdatIT/backend-go/internal/infras/repository/taskgroup"
)

type MockTaskGroupRepository struct {
	mock.Mock
}

func (m *MockTaskGroupRepository) Create(ctx context.Context, item *models.TaskGroup) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockTaskGroupRepository) FindByID(ctx context.Context, id uint64) (*models.TaskGroup, error) {
	args := m.Called(ctx, id)
	var result *models.TaskGroup
	if args.Get(0) != nil {
		result = args.Get(0).(*models.TaskGroup)
	}
	return result, args.Error(1)
}

func (m *MockTaskGroupRepository) FindAllBy(ctx context.Context, params *taskgrouprepo.GetListParams) ([]*models.TaskGroup, int64, error) {
	args := m.Called(ctx, params)
	var results []*models.TaskGroup
	if args.Get(0) != nil {
		results = args.Get(0).([]*models.TaskGroup)
	}
	var total int64
	if args.Get(1) != nil {
		total = args.Get(1).(int64)
	}
	return results, total, args.Error(2)
}

func (m *MockTaskGroupRepository) Update(ctx context.Context, item *models.TaskGroup) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockTaskGroupRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
