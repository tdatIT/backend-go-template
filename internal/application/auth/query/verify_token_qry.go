package query

import (
	"context"
	"errors"
	"log/slog"

	"github.com/tdatIT/backend-go/internal/application/auth/helper"
	"github.com/tdatIT/backend-go/internal/domain/dtos/userdto"
	"github.com/tdatIT/backend-go/internal/infras/repository/session"
	"github.com/tdatIT/backend-go/internal/infras/security"
	"github.com/tdatIT/backend-go/pkgs/decorator"
	"gorm.io/gorm"
)

type IVerifyTokenQuery decorator.QueryHandler[*userdto.VerifyTokenReq, *userdto.VerifyTokenRes]

type verifyTokenQuery struct {
	sessionRepo  session.Repository
	tokenManager security.TokenManager
}

func NewVerifyTokenQuery(sessionRepo session.Repository, tokenManager security.TokenManager) IVerifyTokenQuery {
	return &verifyTokenQuery{
		sessionRepo:  sessionRepo,
		tokenManager: tokenManager,
	}
}

func (v verifyTokenQuery) Handle(ctx context.Context, req *userdto.VerifyTokenReq) (*userdto.VerifyTokenRes, error) {
	claims, err := v.tokenManager.VerifyToken(req.AccessToken)
	if err != nil {
		slog.Warn("failed to verify access token", slog.String("error", err.Error()))
		return nil, helper.ErrInvalidToken
	}

	sessionItem, err := v.sessionRepo.FindByID(ctx, claims.SessionID)
	if err != nil {
		slog.Error("failed to find session by id",
			slog.String("sess_id", claims.SessionID),
			slog.String("error", err.Error()))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrInvalidToken
		}
		return nil, err
	}
	if !sessionItem.IsActive || claims.Subject == "" {
		slog.Warn("invalid session for access token",
			slog.String("sess_id", claims.SessionID),
			slog.Uint64("user_id", sessionItem.UserID))
		return nil, helper.ErrInvalidToken
	}

	return &userdto.VerifyTokenRes{
		Sub:       claims.Subject,
		SessionID: claims.SessionID,
		JTI:       claims.ID,
		ExpiresAt: claims.ExpiresAt.Time.Unix(),
	}, nil
}
