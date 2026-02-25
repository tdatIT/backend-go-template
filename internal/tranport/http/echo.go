package http

import (
	"log/slog"

	"github.com/hellofresh/health-go/v5"
	"github.com/labstack/echo-contrib/v5/echoprometheus"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/tdatIT/backend-go/config"
	"github.com/tdatIT/backend-go/internal/application/auth"
	"github.com/tdatIT/backend-go/internal/tranport/http/handler"
	"github.com/tdatIT/backend-go/internal/tranport/http/helper"
	"github.com/tdatIT/backend-go/internal/tranport/http/router"
	"github.com/tdatIT/backend-go/pkgs/utils/valid"
)

var skipLogURIs = []string{
	"/metrics",
	"/health",
	"/readiness",
	"/liveness",
}

func InitHttpComponent(
	cfg *config.ServiceConfig,
	slogHandler *slog.JSONHandler,
	authApp *auth.Application,
	healthsvc *health.Health,
) *echo.Echo {
	e := echo.New()

	// Recovery middleware to recover from panics and return a 500 error
	e.Use(middleware.Recover())

	// Set up logging middleware
	e.Logger = slog.New(slogHandler)
	e.Use(middleware.RequestLogger())

	// Setup validator
	e.Validator = valid.GetValidator()

	// Global handler error
	e.HTTPErrorHandler = func(c *echo.Context, err error) {
		if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
			if resp.Committed {
				return
			}
		}

		_ = helper.WriteError(c, err)
	}

	// Metrics
	e.Use(echoprometheus.NewMiddleware(cfg.Server.Name)) // adds middleware to gather metrics
	e.GET("/metrics", echoprometheus.NewHandler())       // adds route to serve gathered metrics

	// Health check endpoints
	e.GET("/readiness", echo.WrapHandler(healthsvc.Handler()))
	e.GET("/liveness", func(c *echo.Context) error {
		return c.String(200, "ok")
	})

	// Register routes
	api := e.Group("/api")
	authHandler := handler.NewAuthHandler(authApp)
	router.RegisterAuthRoutes(api, authHandler)

	return e
}
