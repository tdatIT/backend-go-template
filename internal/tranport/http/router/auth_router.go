package router

import (
	"github.com/labstack/echo/v5"
	"github.com/tdatIT/backend-go/internal/tranport/http/handler"
)

func RegisterAuthRoutes(
	router *echo.Group,
	authHandler *handler.AuthHandler,
) {
	auth := router.Group("/v1/auth")
	auth.POST("/login", authHandler.LoginByUserPass)
	auth.POST("/via-google", authHandler.LoginByGoogle)
	auth.POST("/register", authHandler.Register)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.POST("/logout", authHandler.Logout)
}
