package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tdatIT/backend-go/internal/domain/models"
	userrepo "github.com/tdatIT/backend-go/internal/infras/repository/user"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, item *models.User) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint64) (*models.User, error) {
	args := m.Called(ctx, id)
	var result *models.User
	if args.Get(0) != nil {
		result = args.Get(0).(*models.User)
	}
	return result, args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	var result *models.User
	if args.Get(0) != nil {
		result = args.Get(0).(*models.User)
	}
	return result, args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	var result *models.User
	if args.Get(0) != nil {
		result = args.Get(0).(*models.User)
	}
	return result, args.Error(1)
}

func (m *MockUserRepository) FindByOIDC(ctx context.Context, provider string, subject string) (*models.User, error) {
	args := m.Called(ctx, provider, subject)
	var result *models.User
	if args.Get(0) != nil {
		result = args.Get(0).(*models.User)
	}
	return result, args.Error(1)
}

func (m *MockUserRepository) FindAllAndCount(ctx context.Context, params userrepo.GetListParams) ([]*models.User, int64, error) {
	args := m.Called(ctx, params)
	var results []*models.User
	if args.Get(0) != nil {
		results = args.Get(0).([]*models.User)
	}
	var total int64
	if args.Get(1) != nil {
		total = args.Get(1).(int64)
	}
	return results, total, args.Error(2)
}

func (m *MockUserRepository) Update(ctx context.Context, item *models.User) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
