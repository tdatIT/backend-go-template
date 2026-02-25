package query

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tdatIT/backend-go/config"
	"github.com/tdatIT/backend-go/internal/domain/dtos/userdto"
	"github.com/tdatIT/backend-go/internal/domain/models"
	"github.com/tdatIT/backend-go/internal/infras/httpclient/oidc"
	"github.com/tdatIT/backend-go/mocks"
)

func TestLoginByGoogleQuery_Handle_Success(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	sessionRepo := new(mocks.MockSessionRepository)
	tokenManager := new(mocks.MockTokenManager)
	verifier := new(mocks.MockGoogleOIDCVerifier)

	cfg := &config.ServiceConfig{
		Auth: config.Auth{GoogleClientID: "client"},
	}

	info := &oidc.GoogleIDTokenInfo{
		Email:         "user@example.com",
		EmailVerified: "true",
		Subject:       "sub-1",
		Audience:      "client",
		GivenName:     "Test",
		FamilyName:    "User",
	}

	verifier.On("VerifyIDToken", mock.Anything, "idtoken").Return(info, nil)
	userRepo.On("FindByOIDC", mock.Anything, "google", info.Subject).
		Return(&models.User{ID: 5, OidcProvider: "google", OidcSubject: info.Subject}, nil)
	sessionRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Session")).Return(nil)

	accessExp := time.Now().Add(time.Hour)
	tokenManager.On("GenerateTokens", uint64(5), mock.Anything, mock.Anything).
		Return("access", "refresh", accessExp, nil)

	qry := NewLoginByGoogleQuery(userRepo, sessionRepo, tokenManager, verifier, cfg)
	res, err := qry.Handle(context.Background(), &userdto.LoginByGoogleReq{IDToken: "idtoken"})

	require.NoError(t, err)
	require.Equal(t, "access", res.AccessToken)
	require.Equal(t, "refresh", res.RefreshToken)
	require.Equal(t, uint64(5), res.User.ID)

	userRepo.AssertExpectations(t)
	sessionRepo.AssertExpectations(t)
	tokenManager.AssertExpectations(t)
	verifier.AssertExpectations(t)
}
