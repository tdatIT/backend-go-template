package command

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tdatIT/backend-go/internal/application/auth/helper"
	"github.com/tdatIT/backend-go/internal/domain/dtos/userdto"
	"github.com/tdatIT/backend-go/internal/domain/models"
	"github.com/tdatIT/backend-go/internal/infras/security"
	"github.com/tdatIT/backend-go/mocks"
	"gorm.io/gorm"
)

func TestLogoutCommand_Handle_Success(t *testing.T) {
	sessionRepo := new(mocks.MockSessionRepository)
	tokenManager := new(mocks.MockTokenManager)

	claims := &security.CustomClaims{
		SessionID: "sess-1",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
			ID:      "jti-1",
		},
	}

	tokenManager.On("VerifyToken", "access").Return(claims, nil)
	sessionRepo.On("FindBySessionID", mock.Anything, claims.SessionID).
		Return(&models.Session{ID: claims.SessionID, UserID: 1, IsActive: true}, nil)
	sessionRepo.On("Deactivate", mock.Anything, claims.SessionID).Return(nil)

	cmd := NewLogoutCommand(sessionRepo, tokenManager)
	err := cmd.Handle(context.Background(), &userdto.LogoutReq{AccessToken: "access"})

	require.NoError(t, err)
	sessionRepo.AssertExpectations(t)
	tokenManager.AssertExpectations(t)
}

func TestLogoutCommand_Handle_InvalidToken(t *testing.T) {
	sessionRepo := new(mocks.MockSessionRepository)
	tokenManager := new(mocks.MockTokenManager)

	tokenManager.On("VerifyToken", "bad").Return((*security.CustomClaims)(nil), gorm.ErrRecordNotFound)

	cmd := NewLogoutCommand(sessionRepo, tokenManager)
	err := cmd.Handle(context.Background(), &userdto.LogoutReq{AccessToken: "bad"})

	require.ErrorIs(t, err, helper.ErrInvalidToken)
	tokenManager.AssertExpectations(t)
}
