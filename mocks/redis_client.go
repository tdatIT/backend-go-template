package mocks

import (
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Client() redis.UniversalClient {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(redis.UniversalClient)
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}
