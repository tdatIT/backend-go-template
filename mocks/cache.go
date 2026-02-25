package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(ctx context.Context, key string) ([]byte, error) {
	args := m.Called(ctx, key)
	var result []byte
	if args.Get(0) != nil {
		result = args.Get(0).([]byte)
	}
	return result, args.Error(1)
}

func (m *MockCache) Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, val, ttl)
	return args.Error(0)
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	args := m.Called(ctx, key, ttl)
	return args.Error(0)
}
