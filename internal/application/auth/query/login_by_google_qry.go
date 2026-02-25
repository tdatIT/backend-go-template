package query

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tdatIT/backend-go/config"
	"github.com/tdatIT/backend-go/internal/application/auth/helper"
	"github.com/tdatIT/backend-go/internal/domain/dtos/userdto"
	"github.com/tdatIT/backend-go/internal/domain/models"
	"github.com/tdatIT/backend-go/internal/infras/httpclient/oidc"
	"github.com/tdatIT/backend-go/internal/infras/repository/session"
	"github.com/tdatIT/backend-go/internal/infras/repository/user"
	"github.com/tdatIT/backend-go/internal/infras/security"
	"github.com/tdatIT/backend-go/pkgs/decorator"
	"gorm.io/gorm"
)

type ILoginByGoogleQuery decorator.QueryHandler[*userdto.LoginByGoogleReq, *userdto.LoginRes]

type loginByGoogleQuery struct {
	config       *config.ServiceConfig
	userRepo     user.Repository
	sessionRepo  session.Repository
	tokenManager security.TokenManager
	oidcProvider *oidc.GoogleOIDCProvider
}

func NewLoginByGoogleQuery(
	userRepo user.Repository,
	sessionRepo session.Repository,
	tokenManager security.TokenManager,
	oidcProvider *oidc.GoogleOIDCProvider,
	config *config.ServiceConfig,
) ILoginByGoogleQuery {
	return &loginByGoogleQuery{
		config:       config,
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		tokenManager: tokenManager,
		oidcProvider: oidcProvider,
	}
}

func (l loginByGoogleQuery) Handle(ctx context.Context, req *userdto.LoginByGoogleReq) (*userdto.LoginRes, error) {
	info, err := l.oidcProvider.VerifyIDToken(ctx, req.IDToken)
	if err != nil {
		slog.Error("failed to verify google id token", slog.String("error", err.Error()))
		return nil, helper.ErrInvalidToken
	}

	if info.Email == "" || info.EmailVerified != "true" || info.Audience != l.config.Auth.GoogleClientID {
		return nil, helper.ErrInvalidToken
	}

	account, err := l.userRepo.FindByOIDC(ctx, "google", info.Subject)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no user found with OIDC subject, try to find by email to handle case where user previously registered with email/password
			account, err = l.userRepo.FindByEmail(ctx, info.Email)
			if err != nil {
				// If no user found with email, create a new account
				if errors.Is(err, gorm.ErrRecordNotFound) {
					account, err = l.createUserFromGoogle(ctx, info)
					if err != nil {
						return nil, err
					}
				} else {
					slog.Error("failed to find user by email during google login",
						slog.String("email", info.Email),
						slog.String("error", err.Error()))
					return nil, err
				}
			}
		} else {
			slog.Error("failed to find user by OIDC subject during google login",
				slog.String("subject", info.Subject),
				slog.String("error", err.Error()))
			return nil, err
		}
	}

	// If user exists but doesn't have OIDC info, update the account to link it with Google
	if account.OidcProvider == "" || account.OidcSubject == "" {
		account.OidcProvider = "google"
		account.OidcSubject = info.Subject
		if err := l.userRepo.Update(ctx, account); err != nil {
			slog.Error("failed to update user with google OIDC info",
				slog.Uint64("user_id", account.ID),
				slog.String("error", err.Error()))
			return nil, err
		}
	}

	refreshJTI := uuid.NewString()
	now := time.Now()
	sessionItem := &models.Session{
		UserID:     account.ID,
		RefreshJTI: refreshJTI,
		UserAgent:  req.UserAgent,
		IPAddress:  req.IPAddress,
		IsActive:   true,
		LastUsedAt: &now,
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

func (l loginByGoogleQuery) createUserFromGoogle(ctx context.Context, info *oidc.GoogleIDTokenInfo) (*models.User, error) {
	firstName := strings.TrimSpace(info.GivenName)
	lastName := strings.TrimSpace(info.FamilyName)
	if firstName == "" {
		firstName = strings.TrimSpace(info.Name)
	}
	if firstName == "" {
		firstName = "Google"
	}
	if lastName == "" {
		lastName = "User"
	}

	password, err := security.HashPassword(uuid.NewString(), security.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password for new google user",
			slog.String("email", info.Email),
			slog.String("error", err.Error()))
		return nil, err
	}

	item := &models.User{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        info.Email,
		Username:     info.Email,
		PasswordHash: password,
		IsActive:     true,
		OidcProvider: "google",
		OidcSubject:  info.Subject,
	}

	if err := l.userRepo.Create(ctx, item); err != nil {
		slog.Error("failed to create user from google info",
			slog.String("email", info.Email),
			slog.String("error", err.Error()))
		return nil, err
	}

	return item, nil
}
