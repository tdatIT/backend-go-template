package query

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tdatIT/backend-go/internal/application/auth/helper"
	"github.com/tdatIT/backend-go/internal/domain/dtos/userdto"
	"github.com/tdatIT/backend-go/internal/domain/models"
	"github.com/tdatIT/backend-go/internal/infras/oidc/google"
	"github.com/tdatIT/backend-go/internal/infras/repository/user"
	"github.com/tdatIT/backend-go/internal/infras/security"
	"github.com/tdatIT/backend-go/pkgs/decorator"
	"gorm.io/gorm"
)

type ILoginByGoogleQuery decorator.QueryHandler[*userdto.LoginByGoogleReq, *userdto.LoginRes]

type loginByGoogleQuery struct {
	userRepo       user.Repository
	tokenManager   security.TokenManager
	googleVerifier *google.TokenVerifier
	googleClientID string
}

func NewLoginByGoogleQuery(
	userRepo user.Repository,
	tokenManager security.TokenManager,
	googleVerifier *google.TokenVerifier,
	googleClientID string,
) ILoginByGoogleQuery {
	return &loginByGoogleQuery{
		userRepo:       userRepo,
		tokenManager:   tokenManager,
		googleVerifier: googleVerifier,
		googleClientID: googleClientID,
	}
}

func (l loginByGoogleQuery) Handle(ctx context.Context, req *userdto.LoginByGoogleReq) (*userdto.LoginRes, error) {
	if req == nil || strings.TrimSpace(req.IDToken) == "" {
		return nil, errors.New("id token is required")
	}

	info, err := l.googleVerifier.VerifyIDToken(ctx, req.IDToken)
	if err != nil {
		slog.Error("failed to verify google id token", slog.String("error", err.Error()))
		return nil, helper.ErrInvalidToken
	}

	if info.Email == "" {
		return nil, helper.ErrInvalidToken
	}
	if info.EmailVerified != "true" {
		return nil, helper.ErrInvalidToken
	}
	if l.googleClientID != "" && info.Audience != l.googleClientID {
		return nil, helper.ErrInvalidToken
	}

	account, err := l.userRepo.FindByOIDC(ctx, "google", info.Subject)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			account, err = l.userRepo.FindByEmail(ctx, info.Email)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					account, err = l.createUserFromGoogle(ctx, info)
					if err != nil {
						return nil, err
					}
				} else {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	}

	if account.OidcProvider == "" || account.OidcSubject == "" {
		account.OidcProvider = "google"
		account.OidcSubject = info.Subject
		if err := l.userRepo.Update(ctx, account); err != nil {
			return nil, err
		}
	}

	accessToken, refreshToken, accessExp, err := l.tokenManager.GenerateTokens(account.ID)
	if err != nil {
		return nil, err
	}

	refreshHash, err := security.HashPassword(refreshToken, security.DefaultCost)
	if err != nil {
		return nil, err
	}

	account.RefreshTokenHash = refreshHash
	if err := l.userRepo.Update(ctx, account); err != nil {
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

func (l loginByGoogleQuery) createUserFromGoogle(ctx context.Context, info *google.IDTokenInfo) (*models.User, error) {
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
		return nil, err
	}

	return item, nil
}
