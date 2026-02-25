package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tdatIT/backend-go/internal/domain/models"
)

type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Create(ctx context.Context, item *models.Session) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockSessionRepository) FindByID(ctx context.Context, id string) (*models.Session, error) {
	args := m.Called(ctx, id)
	var result *models.Session
	if args.Get(0) != nil {
		result = args.Get(0).(*models.Session)
	}
	return result, args.Error(1)
}

func (m *MockSessionRepository) FindByRefreshJTI(ctx context.Context, jti string) (*models.Session, error) {
	args := m.Called(ctx, jti)
	var result *models.Session
	if args.Get(0) != nil {
		result = args.Get(0).(*models.Session)
	}
	return result, args.Error(1)
}

func (m *MockSessionRepository) FindBySessionID(ctx context.Context, sessionID string) (*models.Session, error) {
	args := m.Called(ctx, sessionID)
	var result *models.Session
	if args.Get(0) != nil {
		result = args.Get(0).(*models.Session)
	}
	return result, args.Error(1)
}

func (m *MockSessionRepository) Update(ctx context.Context, item *models.Session) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockSessionRepository) RotateRefreshJTI(ctx context.Context, id string, oldJTI string, newJTI string) error {
	args := m.Called(ctx, id, oldJTI, newJTI)
	return args.Error(0)
}

func (m *MockSessionRepository) Deactivate(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
