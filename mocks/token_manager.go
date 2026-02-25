package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tdatIT/backend-go/internal/infras/security"
)

type MockTokenManager struct {
	mock.Mock
}

func (m *MockTokenManager) GenerateTokens(subjectID uint64, sessionID string, refreshJTI string) (string, string, time.Time, error) {
	args := m.Called(subjectID, sessionID, refreshJTI)
	return args.String(0), args.String(1), args.Get(2).(time.Time), args.Error(3)
}

func (m *MockTokenManager) VerifyToken(tokenString string) (*security.CustomClaims, error) {
	args := m.Called(tokenString)
	var result *security.CustomClaims
	if args.Get(0) != nil {
		result = args.Get(0).(*security.CustomClaims)
	}
	return result, args.Error(1)
}
