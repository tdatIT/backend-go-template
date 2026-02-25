package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tdatIT/backend-go/internal/infras/httpclient/oidc"
)

type MockGoogleOIDCVerifier struct {
	mock.Mock
}

func (m *MockGoogleOIDCVerifier) VerifyIDToken(ctx context.Context, idToken string) (*oidc.GoogleIDTokenInfo, error) {
	args := m.Called(ctx, idToken)
	var result *oidc.GoogleIDTokenInfo
	if args.Get(0) != nil {
		result = args.Get(0).(*oidc.GoogleIDTokenInfo)
	}
	return result, args.Error(1)
}
