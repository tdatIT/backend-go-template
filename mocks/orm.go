package mocks

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockORM struct {
	mock.Mock
}

func (m *MockORM) GormDB() *gorm.DB {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*gorm.DB)
}

func (m *MockORM) SqlDB() *sql.DB {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*sql.DB)
}

func (m *MockORM) Close() error {
	args := m.Called()
	return args.Error(0)
}
