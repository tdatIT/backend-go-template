package query

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/tdatIT/backend-go/internal/application/auth/helper"
	"github.com/tdatIT/backend-go/internal/domain/dtos/userdto"
	"github.com/tdatIT/backend-go/internal/infras/repository/user"
	"github.com/tdatIT/backend-go/internal/infras/security"
	"github.com/tdatIT/backend-go/pkgs/decorator"
)

type ILoginByUsrnameAndPwdQuery decorator.QueryHandler[*userdto.LoginByUserPassReq, *userdto.LoginRes]

type loginByUsrnameAndPwdQuery struct {
	userRepo     user.Repository
	tokenManager security.TokenManager
}

func NewLoginByUsrnameAndPwdQuery(
	userRepo user.Repository,
	tokenManager security.TokenManager,
) ILoginByUsrnameAndPwdQuery {
	return &loginByUsrnameAndPwdQuery{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}

func (l loginByUsrnameAndPwdQuery) Handle(ctx context.Context, req *userdto.LoginByUserPassReq) (*userdto.LoginRes, error) {
	if req == nil {
		return nil, errors.New("login request is required")
	}

	account, err := l.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
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

	accessToken, refreshToken, accessExp, err := l.tokenManager.GenerateTokens(account.ID)
	if err != nil {
		slog.Error("failed to generate tokens for user",
			slog.Uint64("userID", account.ID),
			slog.String("error", err.Error()))
		return nil, err
	}

	refreshHash, err := security.HashPassword(refreshToken, security.DefaultCost)
	if err != nil {
		slog.Error("failed to hash refresh token",
			slog.Uint64("userID", account.ID),
			slog.String("error", err.Error()))
		return nil, err
	}

	account.RefreshTokenHash = refreshHash
	if err := l.userRepo.Update(ctx, account); err != nil {
		slog.Error("failed to update refresh token for user",
			slog.Uint64("userID", account.ID),
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
