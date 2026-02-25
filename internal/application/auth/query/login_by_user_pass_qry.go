package query

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/tdatIT/backend-go/internal/application/auth/helper"
	"github.com/tdatIT/backend-go/internal/domain/dtos/userdto"
	"github.com/tdatIT/backend-go/internal/domain/models"
	"github.com/tdatIT/backend-go/internal/infras/repository/session"
	"github.com/tdatIT/backend-go/internal/infras/repository/user"
	"github.com/tdatIT/backend-go/internal/infras/security"
	"github.com/tdatIT/backend-go/pkgs/decorator"
	"github.com/tdatIT/backend-go/pkgs/utils/genid"
	"gorm.io/gorm"
)

type ILoginByUsrnameAndPwdQuery decorator.QueryHandler[*userdto.LoginByUserPassReq, *userdto.LoginRes]

type loginByUsrnameAndPwdQuery struct {
	userRepo     user.Repository
	sessionRepo  session.Repository
	tokenManager security.TokenManager
}

func NewLoginByUsrnameAndPwdQuery(
	userRepo user.Repository,
	sessionRepo session.Repository,
	tokenManager security.TokenManager,
) ILoginByUsrnameAndPwdQuery {
	return &loginByUsrnameAndPwdQuery{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		tokenManager: tokenManager,
	}
}

func (l loginByUsrnameAndPwdQuery) Handle(ctx context.Context, req *userdto.LoginByUserPassReq) (*userdto.LoginRes, error) {
	account, err := l.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrUserNotFound
		}

		slog.Error("failed to find user by username",
			slog.String("username", req.Username),
			slog.String("error", err.Error()))
		return nil, err
	}

	if err := security.ComparePassword(account.PasswordHash, req.Password); err != nil {
		slog.Warn("invalid password for user",
			slog.String("username", req.Username),
			slog.String("error", err.Error()))
		return nil, helper.ErrInvalidUserOrPwd
	}

	refreshJTI := uuid.NewString()
	sessionItem := &models.Session{
		ID:         genid.GenerateNanoID(),
		UserID:     account.ID,
		RefreshJTI: refreshJTI,
		UserAgent:  req.UserAgent,
		IPAddress:  req.IPAddress,
		IsActive:   true,
		LastUsedAt: new(time.Now()),
	}

	if err := l.sessionRepo.Create(ctx, sessionItem); err != nil {
		slog.Error("failed to create session",
			slog.Uint64("user_id", account.ID),
			slog.String("error", err.Error()))
		return nil, err
	}

	accessToken, refreshToken, accessExp, err := l.tokenManager.GenerateTokens(account.ID, sessionItem.ID, refreshJTI)
	if err != nil {
		slog.Error("failed to generate tokens for user",
			slog.Uint64("user_id", account.ID),
			slog.String("error", err.Error()))
		return nil, err
	}

	return &userdto.LoginRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(time.Until(accessExp).Seconds()),
		User: &userdto.UserProfileRes{
			ID:        account.ID,
			FirstName: account.FirstName,
			LastName:  account.LastName,
			Email:     account.Email,
			Username:  account.Username,
		},
	}, nil
}
