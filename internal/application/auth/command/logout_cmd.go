package command

import (
	"context"
	"errors"
	"log/slog"

	"github.com/tdatIT/backend-go/internal/application/auth/helper"
	"github.com/tdatIT/backend-go/internal/domain/dtos/userdto"
	"github.com/tdatIT/backend-go/internal/infras/repository/user"
	"github.com/tdatIT/backend-go/internal/infras/security"
	"github.com/tdatIT/backend-go/pkgs/decorator"
	"gorm.io/gorm"
)

type ILogoutCommand decorator.CommandHandler[*userdto.LogoutReq]

type logoutCommand struct {
	userRepo user.Repository
}

func NewLogoutCommand(userRepo user.Repository) ILogoutCommand {
	return &logoutCommand{userRepo: userRepo}
}

func (l logoutCommand) Handle(ctx context.Context, req *userdto.LogoutReq) error {
	account, err := l.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		slog.Error("failed to find user by id", slog.Uint64("user_id", req.UserID), slog.String("error", err.Error()))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.ErrInvalidToken
		}
		return err
	}

	if account.RefreshTokenHash == "" {
		return helper.ErrInvalidToken
	}
	if err := security.ComparePassword(account.RefreshTokenHash, req.RefreshToken); err != nil {
		return helper.ErrInvalidToken
	}

	account.RefreshTokenHash = ""
	return l.userRepo.Update(ctx, account)
}
