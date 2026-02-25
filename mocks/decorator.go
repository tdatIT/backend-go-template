package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockCommandHandler[T any] struct {
	mock.Mock
}

func (m *MockCommandHandler[T]) Handle(ctx context.Context, req T) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

type MockCommandReturnHandler[T any, E any] struct {
	mock.Mock
}

func (m *MockCommandReturnHandler[T, E]) Handle(ctx context.Context, req T) (E, error) {
	args := m.Called(ctx, req)
	var result E
	if args.Get(0) != nil {
		result = args.Get(0).(E)
	}
	return result, args.Error(1)
}

type MockQueryHandler[T any, E any] struct {
	mock.Mock
}

func (m *MockQueryHandler[T, E]) Handle(ctx context.Context, req T) (E, error) {
	args := m.Called(ctx, req)
	var result E
	if args.Get(0) != nil {
		result = args.Get(0).(E)
	}
	return result, args.Error(1)
}
