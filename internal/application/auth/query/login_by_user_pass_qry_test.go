package query

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tdatIT/backend-go/internal/application/auth/helper"
	"github.com/tdatIT/backend-go/internal/domain/dtos/userdto"
	"github.com/tdatIT/backend-go/internal/domain/models"
	"github.com/tdatIT/backend-go/internal/infras/security"
	"github.com/tdatIT/backend-go/mocks"
	"gorm.io/gorm"
)

func TestLoginByUsrnameAndPwdQuery_Handle_Success(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	sessionRepo := new(mocks.MockSessionRepository)
	tokenManager := new(mocks.MockTokenManager)

	passwordHash, err := security.HashPassword("pass1234", security.DefaultCost)
	require.NoError(t, err)

	account := &models.User{ID: 10, Username: "user", PasswordHash: passwordHash}
	userRepo.On("FindByUsername", mock.Anything, "user").Return(account, nil)
	sessionRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Session")).Return(nil)

	accessExp := time.Now().Add(time.Hour)
	tokenManager.On("GenerateTokens", uint64(10), mock.Anything, mock.Anything).
		Return("access", "refresh", accessExp, nil)

	qry := NewLoginByUsrnameAndPwdQuery(userRepo, sessionRepo, tokenManager)
	res, err := qry.Handle(context.Background(), &userdto.LoginByUserPassReq{
		Username:  "user",
		Password:  "pass1234",
		UserAgent: "ua",
		IPAddress: "127.0.0.1",
	})

	require.NoError(t, err)
	require.Equal(t, "access", res.AccessToken)
	require.Equal(t, "refresh", res.RefreshToken)

	userRepo.AssertExpectations(t)
	sessionRepo.AssertExpectations(t)
	tokenManager.AssertExpectations(t)
}

func TestLoginByUsrnameAndPwdQuery_Handle_UserNotFound(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	sessionRepo := new(mocks.MockSessionRepository)
	tokenManager := new(mocks.MockTokenManager)

	userRepo.On("FindByUsername", mock.Anything, "user").
		Return((*models.User)(nil), gorm.ErrRecordNotFound)

	qry := NewLoginByUsrnameAndPwdQuery(userRepo, sessionRepo, tokenManager)
	res, err := qry.Handle(context.Background(), &userdto.LoginByUserPassReq{
		Username: "user",
		Password: "pass1234",
	})

	require.Nil(t, res)
	require.ErrorIs(t, err, helper.ErrUserNotFound)

	userRepo.AssertExpectations(t)
}
