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

func TestVerifyTokenQuery_Handle_Success(t *testing.T) {
	sessionRepo := new(mocks.MockSessionRepository)
	tokenManager := new(mocks.MockTokenManager)

	claims := &security.CustomClaims{
		SessionID: "sess-1",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			ID:        "access-1",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	sessionRepo.On("FindByID", mock.Anything, claims.SessionID).
		Return(&models.Session{ID: claims.SessionID, UserID: 1, IsActive: true}, nil)
	tokenManager.On("VerifyToken", "access").Return(claims, nil)

	qry := NewVerifyTokenQuery(sessionRepo, tokenManager)
	res, err := qry.Handle(context.Background(), &userdto.VerifyTokenReq{AccessToken: "access"})

	require.NoError(t, err)
	require.Equal(t, claims.Subject, res.Sub)
	require.Equal(t, claims.SessionID, res.SessionID)
	require.Equal(t, claims.ID, res.JTI)

	sessionRepo.AssertExpectations(t)
	tokenManager.AssertExpectations(t)
}

func TestVerifyTokenQuery_Handle_InvalidToken(t *testing.T) {
	sessionRepo := new(mocks.MockSessionRepository)
	tokenManager := new(mocks.MockTokenManager)

	tokenManager.On("VerifyToken", "bad").Return((*security.CustomClaims)(nil), gorm.ErrRecordNotFound)

	qry := NewVerifyTokenQuery(sessionRepo, tokenManager)
	res, err := qry.Handle(context.Background(), &userdto.VerifyTokenReq{AccessToken: "bad"})

	require.Nil(t, res)
	require.ErrorIs(t, err, helper.ErrInvalidToken)

	tokenManager.AssertExpectations(t)
}
