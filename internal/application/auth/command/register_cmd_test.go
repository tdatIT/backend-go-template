package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tdatIT/backend-go/internal/application/auth/helper"
	"github.com/tdatIT/backend-go/internal/domain/dtos/userdto"
	"github.com/tdatIT/backend-go/internal/domain/models"
	"github.com/tdatIT/backend-go/mocks"
	"gorm.io/gorm"
)

func TestRegisterCommand_Handle_Success(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	sessionRepo := new(mocks.MockSessionRepository)
	tokenManager := new(mocks.MockTokenManager)

	req := &userdto.RegisterReq{
		Username:  "test@example.com",
		Password:  "pass1234",
		FirstName: "Test",
		LastName:  "User",
		UserAgent: "ua",
		IPAddress: "127.0.0.1",
	}

	userRepo.On("FindByUsername", mock.Anything, req.Username).
		Return((*models.User)(nil), gorm.ErrRecordNotFound)
	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).
		Run(func(args mock.Arguments) {
			item := args.Get(1).(*models.User)
			item.ID = 1
		}).Return(nil)

	sessionRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Session")).Return(nil)

	accessExp := time.Now().Add(time.Hour)
	tokenManager.On("GenerateTokens", uint64(1), mock.Anything, mock.Anything).
		Return("access", "refresh", accessExp, nil)

	cmd := NewRegisterCommand(userRepo, sessionRepo, tokenManager)
	res, err := cmd.Handle(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, "access", res.AccessToken)
	require.Equal(t, "refresh", res.RefreshToken)
	require.Equal(t, uint64(1), res.User.ID)

	userRepo.AssertExpectations(t)
	sessionRepo.AssertExpectations(t)
	tokenManager.AssertExpectations(t)
}

func TestRegisterCommand_Handle_UserAlreadyExists(t *testing.T) {
	userRepo := new(mocks.MockUserRepository)
	sessionRepo := new(mocks.MockSessionRepository)
	tokenManager := new(mocks.MockTokenManager)

	req := &userdto.RegisterReq{
		Username:  "test@example.com",
		Password:  "pass1234",
		FirstName: "Test",
	}

	userRepo.On("FindByUsername", mock.Anything, req.Username).
		Return(&models.User{ID: 9}, nil)

	cmd := NewRegisterCommand(userRepo, sessionRepo, tokenManager)
	res, err := cmd.Handle(context.Background(), req)

	require.Nil(t, res)
	require.ErrorIs(t, err, helper.ErrUserAlreadyExists)

	userRepo.AssertExpectations(t)
}

func TestRegisterCommand_Handle_InvalidRequest(t *testing.T) {
	cmd := NewRegisterCommand(nil, nil, nil)
	res, err := cmd.Handle(context.Background(), nil)

	require.Nil(t, res)
	require.Error(t, err)
}

func TestRegisterCommand_Handle_InvalidPassword(t *testing.T) {
	req := &userdto.RegisterReq{
		Username:  "test@example.com",
		Password:  "",
		FirstName: "Test",
	}

	cmd := NewRegisterCommand(nil, nil, nil)
	res, err := cmd.Handle(context.Background(), req)

	require.Nil(t, res)
	require.Error(t, err)
}
