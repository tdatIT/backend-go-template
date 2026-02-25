package command

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

type ILogoutCommand decorator.CommandHandler[*userdto.LogoutReq]

type logoutCommand struct {
	sessionRepo  session.Repository
	tokenManager security.TokenManager
}

func NewLogoutCommand(sessionRepo session.Repository, tokenManager security.TokenManager) ILogoutCommand {
	return &logoutCommand{
		sessionRepo:  sessionRepo,
		tokenManager: tokenManager,
	}
}

func (l logoutCommand) Handle(ctx context.Context, req *userdto.LogoutReq) error {
	claims, err := l.tokenManager.VerifyToken(req.AccessToken)
	if err != nil {
		slog.Warn("failed to verify access token", slog.String("error", err.Error()))
		return helper.ErrInvalidToken
	}

	sessionItem, err := l.sessionRepo.FindBySessionID(ctx, claims.SessionID)
	if err != nil {
		slog.Warn("failed to find session by refresh jti",
			slog.String("sess_id", claims.SessionID),
			slog.String("error", err.Error()))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.ErrInvalidToken
		}
		return err
	}

	if sessionItem.ID != claims.SessionID || claims.Subject == "" || !sessionItem.IsActive {
		slog.Warn("invalid session for refresh token",
			slog.String("sess_id", claims.SessionID),
			slog.Uint64("user_id", sessionItem.UserID),
			slog.Bool("is_active", sessionItem.IsActive))
		return helper.ErrInvalidToken
	}

	if err := l.sessionRepo.Deactivate(ctx, sessionItem.ID); err != nil {
		slog.Error("failed to deactivate session",
			slog.String("sess_id", sessionItem.ID),
			slog.String("error", err.Error()))
		return err
	}

	return nil
}
