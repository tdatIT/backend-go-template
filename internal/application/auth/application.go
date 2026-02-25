package auth

import (
	"github.com/tdatIT/backend-go/config"
	"github.com/tdatIT/backend-go/internal/application/auth/command"
	"github.com/tdatIT/backend-go/internal/application/auth/query"
	"github.com/tdatIT/backend-go/internal/infras/httpclient/oidc"
	"github.com/tdatIT/backend-go/internal/infras/repository/session"
	"github.com/tdatIT/backend-go/internal/infras/repository/user"
	"github.com/tdatIT/backend-go/internal/infras/security"
)

type queries struct {
	LoginByUsernameAndPassword query.ILoginByUsrnameAndPwdQuery
	LoginByGoogle              query.ILoginByGoogleQuery
	RefreshToken               query.IRefreshTokenQuery
	VerifyToken                query.IVerifyTokenQuery
}

type commands struct {
	Register command.IRegisterCommand
	Logout   command.ILogoutCommand
}

type Application struct {
	Queries  *queries
	Commands *commands
}

func NewApplication(
	config *config.ServiceConfig,
	userRepo user.Repository,
	sessionRepo session.Repository,
	tokenManager security.TokenManager,
) *Application {
	// Initialize essential components
	googleOIDC := oidc.NewGoogleOIDCProvider(config)

	return &Application{
		Queries: &queries{
			LoginByUsernameAndPassword: query.NewLoginByUsrnameAndPwdQuery(userRepo, sessionRepo, tokenManager),
			LoginByGoogle:              query.NewLoginByGoogleQuery(userRepo, sessionRepo, tokenManager, googleOIDC, config),
			RefreshToken:               query.NewRefreshTokenQuery(sessionRepo, tokenManager),
			VerifyToken:                query.NewVerifyTokenQuery(sessionRepo, tokenManager),
		},
		Commands: &commands{
			Register: command.NewRegisterCommand(userRepo, sessionRepo, tokenManager),
			Logout:   command.NewLogoutCommand(sessionRepo, tokenManager),
		},
	}
}
