package server

import (
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/tdatIT/backend-go/config"
	"github.com/tdatIT/backend-go/internal/application/auth"
	"github.com/tdatIT/backend-go/internal/infras/repository/session"
	"github.com/tdatIT/backend-go/internal/infras/repository/user"
	"github.com/tdatIT/backend-go/internal/infras/security"
	httpComponent "github.com/tdatIT/backend-go/internal/tranport/http"
	"github.com/tdatIT/backend-go/pkgs/cache"
	"github.com/tdatIT/backend-go/pkgs/db/orm"
	"github.com/tdatIT/backend-go/pkgs/db/rdclient"
	htlcheck "github.com/tdatIT/backend-go/pkgs/hltcheck"
	"github.com/tdatIT/backend-go/pkgs/logger"
)

type Service struct {
	_config *config.ServiceConfig
	_echo   *echo.Echo
}

func InitServer() (*Service, error) {
	svcConfig, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to load config", slog.String("error", err.Error()))
		return nil, err
	}

	jsonLogHandler := logger.NewJsonSlogHandler(&logger.SlogConfig{
		Level:       svcConfig.Logger.Level,
		ServiceName: svcConfig.Server.Name,
	})
	slog.SetDefault(slog.New(jsonLogHandler))

	//connect database, cache
	database := orm.NewDBConnection(svcConfig)
	redis := rdclient.NewRedisClient(svcConfig)
	cacheEngine := cache.NewRedisCacheClient(redis)

	//init repositories, services, handlers
	userRepo := user.NewRepository(database)
	sessRepo := session.NewRepository(database, cacheEngine, svcConfig.Auth.RefreshTokenTTL)

	tokenManager := security.NewJWTTokenManager(security.JWTConfig{
		Secret:          svcConfig.Auth.JWTSecret,
		AccessTokenTTL:  svcConfig.Auth.AccessTokenTTL,
		RefreshTokenTTL: svcConfig.Auth.RefreshTokenTTL,
	})

	authApp := auth.NewApplication(svcConfig, userRepo, sessRepo, tokenManager)

	//health service
	healthsvc, _ := htlcheck.NewHealthCheckService(svcConfig, database, redis)

	//init http server
	echoHttp := httpComponent.InitHttpComponent(svcConfig, jsonLogHandler, authApp, healthsvc)

	return &Service{
		_config: svcConfig,
		_echo:   echoHttp,
	}, nil
}

func (s *Service) StartHTTP() error {
	slog.Info("starting HTTP server", slog.String("address", s._config.Server.HttpPort))
	return s._echo.Start(s._config.Server.HttpPort)
}
