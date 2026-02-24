package auth

import (
	"github.com/tdatIT/backend-go/internal/application/auth/command"
	"github.com/tdatIT/backend-go/internal/application/auth/query"
	"github.com/tdatIT/backend-go/internal/infras/oidc/google"
	"github.com/tdatIT/backend-go/internal/infras/repository/user"
	"github.com/tdatIT/backend-go/internal/infras/security"
)

type queries struct {
	LoginByUsernameAndPassword query.ILoginByUsrnameAndPwdQuery
	LoginByGoogle              query.ILoginByGoogleQuery
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
	userRepo user.Repository,
	tokenManager security.TokenManager,
	googleClientID string,
) *Application {
	googleVerifier := google.NewTokenVerifier()

	return &Application{
		Queries: &queries{
			LoginByUsernameAndPassword: query.NewLoginByUsrnameAndPwdQuery(userRepo, tokenManager),
			LoginByGoogle:              query.NewLoginByGoogleQuery(userRepo, tokenManager, googleVerifier, googleClientID),
			VerifyToken:                query.NewVerifyTokenQuery(tokenManager),
		},
		Commands: &commands{
			Register: command.NewRegisterCommand(userRepo, tokenManager),
			Logout:   command.NewLogoutCommand(userRepo),
		},
	}
}
