package command

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tdatIT/backend-go/internal/application/auth/helper"
	"github.com/tdatIT/backend-go/internal/domain/dtos/userdto"
	"github.com/tdatIT/backend-go/internal/domain/models"
	"github.com/tdatIT/backend-go/internal/infras/repository/session"
	"github.com/tdatIT/backend-go/internal/infras/repository/user"
	"github.com/tdatIT/backend-go/internal/infras/security"
	"github.com/tdatIT/backend-go/pkgs/decorator"
	"gorm.io/gorm"
)

type IRegisterCommand decorator.CommandReturnHandler[*userdto.RegisterReq, *userdto.LoginRes]

type registerCommand struct {
	userRepo     user.Repository
	sessionRepo  session.Repository
	tokenManager security.TokenManager
}

func NewRegisterCommand(userRepo user.Repository, sessionRepo session.Repository, tokenManager security.TokenManager) IRegisterCommand {
	return &registerCommand{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		tokenManager: tokenManager,
	}
}

func (r registerCommand) Handle(ctx context.Context, req *userdto.RegisterReq) (*userdto.LoginRes, error) {
	if req == nil {
		return nil, errors.New("register request is required")
	}
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" || strings.TrimSpace(req.FirstName) == "" {
		return nil, errors.New("username, password, and first name are required")
	}

	if _, err := r.userRepo.FindByUsername(ctx, req.Username); err == nil {
		return nil, helper.ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	lastName := strings.TrimSpace(req.LastName)
	if lastName == "" {
		lastName = "User"
	}

	passwordHash, err := security.HashPassword(req.Password, security.DefaultCost)
	if err != nil {
		return nil, err
	}

	item := &models.User{
		FirstName:    strings.TrimSpace(req.FirstName),
		LastName:     lastName,
		Email:        strings.TrimSpace(req.Username),
		Username:     strings.TrimSpace(req.Username),
		PasswordHash: passwordHash,
		IsActive:     true,
	}

	if err := r.userRepo.Create(ctx, item); err != nil {
		return nil, err
	}

	refreshJTI := uuid.NewString()
	now := time.Now()
	sessionItem := &models.Session{
		UserID:     item.ID,
		RefreshJTI: refreshJTI,
		UserAgent:  req.UserAgent,
		IPAddress:  req.IPAddress,
		IsActive:   true,
		LastUsedAt: &now,
	}
	if err := r.sessionRepo.Create(ctx, sessionItem); err != nil {
		return nil, err
	}

	accessToken, refreshToken, accessExp, err := r.tokenManager.GenerateTokens(item.ID, sessionItem.ID, refreshJTI)
	if err != nil {
		return nil, err
	}

	return &userdto.LoginRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(time.Until(accessExp).Seconds()),
		User: &userdto.UserProfileRes{
			ID:        item.ID,
			FirstName: item.FirstName,
			LastName:  item.LastName,
			Email:     item.Email,
			Username:  item.Username,
		},
	}, nil
}
