package query

import (
	"context"
	"testing"
	"time"

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

func TestRefreshTokenQuery_Handle_Success(t *testing.T) {
	sessionRepo := new(mocks.MockSessionRepository)
	tokenManager := new(mocks.MockTokenManager)

	claims := &security.CustomClaims{
		SessionID: "sess-1",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
			ID:      "refresh-1",
		},
	}

	sessionItem := &models.Session{
		ID:         "sess-1",
		UserID:     1,
		RefreshJTI: "refresh-1",
		IsActive:   true,
	}

	tokenManager.On("VerifyToken", "refresh").Return(claims, nil)
	sessionRepo.On("FindBySessionID", mock.Anything, claims.SessionID).Return(sessionItem, nil)

	accessExp := time.Now().Add(time.Hour)
	tokenManager.On("GenerateTokens", uint64(1), sessionItem.ID, mock.Anything).
		Return("access", "refresh-new", accessExp, nil)

	sessionRepo.On("RotateRefreshJTI", mock.Anything, sessionItem.ID, sessionItem.RefreshJTI, mock.Anything).Return(nil)

	qry := NewRefreshTokenQuery(sessionRepo, tokenManager)
	res, err := qry.Handle(context.Background(), &userdto.RefreshTokenReq{RefreshToken: "refresh"})

	require.NoError(t, err)
	require.Equal(t, "access", res.AccessToken)
	require.Equal(t, "refresh-new", res.RefreshToken)

	sessionRepo.AssertExpectations(t)
	tokenManager.AssertExpectations(t)
}

func TestRefreshTokenQuery_Handle_InvalidToken(t *testing.T) {
	sessionRepo := new(mocks.MockSessionRepository)
	tokenManager := new(mocks.MockTokenManager)

	tokenManager.On("VerifyToken", "bad").Return((*security.CustomClaims)(nil), gorm.ErrRecordNotFound)

	qry := NewRefreshTokenQuery(sessionRepo, tokenManager)
	res, err := qry.Handle(context.Background(), &userdto.RefreshTokenReq{RefreshToken: "bad"})

	require.Nil(t, res)
	require.ErrorIs(t, err, helper.ErrInvalidToken)

	tokenManager.AssertExpectations(t)
}
