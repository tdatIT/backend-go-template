package query

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/tdatIT/backend-go/internal/application/auth/helper"
	"github.com/tdatIT/backend-go/internal/domain/dtos/userdto"
	"github.com/tdatIT/backend-go/internal/infras/repository/session"
	"github.com/tdatIT/backend-go/internal/infras/security"
	"github.com/tdatIT/backend-go/pkgs/decorator"
	"gorm.io/gorm"
)

type IRefreshTokenQuery decorator.QueryHandler[*userdto.RefreshTokenReq, *userdto.RefreshTokenRes]

type refreshTokenQuery struct {
	sessionRepo  session.Repository
	tokenManager security.TokenManager
}

func NewRefreshTokenQuery(sessionRepo session.Repository, tokenManager security.TokenManager) IRefreshTokenQuery {
	return &refreshTokenQuery{
		sessionRepo:  sessionRepo,
		tokenManager: tokenManager,
	}
}

func (r refreshTokenQuery) Handle(ctx context.Context, req *userdto.RefreshTokenReq) (*userdto.RefreshTokenRes, error) {
	claims, err := r.tokenManager.VerifyToken(req.RefreshToken)
	if err != nil {
		slog.Warn("failed to verify refresh token", slog.String("error", err.Error()))
		return nil, helper.ErrInvalidToken
	}

	sessionItem, err := r.sessionRepo.FindBySessionID(ctx, claims.SessionID)
	if err != nil {
		slog.Warn("failed to find session by refresh JTI",
			slog.String("refresh_jti", claims.ID),
			slog.String("sess_id", claims.SessionID),
			slog.String("error", err.Error()))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrInvalidToken
		}
		return nil, err
	}

	if !sessionItem.IsActive || sessionItem.ID != claims.SessionID ||
		claims.Subject == "" || sessionItem.RefreshJTI != claims.ID {
		slog.Warn("invalid session for refresh token",
			slog.String("sess_id", sessionItem.ID),
			slog.Uint64("user_id", sessionItem.UserID),
			slog.String("token_sub", claims.Subject),
			slog.String("jti", sessionItem.RefreshJTI),
			slog.String("error", "session is not active or does not match token claims"))
		return nil, helper.ErrInvalidToken
	}

	newRefreshJTI := uuid.NewString()
	accessToken, refreshToken, accessExp, err := r.tokenManager.GenerateTokens(sessionItem.UserID, sessionItem.ID, newRefreshJTI)
	if err != nil {
		slog.Error("failed to generate new tokens",
			slog.Uint64("user_id", sessionItem.UserID),
			slog.String("sess_id", sessionItem.ID),
			slog.String("error", err.Error()))
		return nil, err
	}

	if err := r.sessionRepo.RotateRefreshJTI(ctx, sessionItem.ID, sessionItem.RefreshJTI, newRefreshJTI); err != nil {
		slog.Error("failed to rotate refresh JTI",
			slog.String("sess_id", sessionItem.ID),
			slog.String("old_jti", sessionItem.RefreshJTI),
			slog.String("new_jti", newRefreshJTI),
			slog.String("error", err.Error()))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrInvalidToken
		}
		return nil, err
	}

	return &userdto.RefreshTokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(time.Until(accessExp).Seconds()),
	}, nil
}
