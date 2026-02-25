package handler

import (
	"log/slog"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/tdatIT/backend-go/internal/application/auth"
	"github.com/tdatIT/backend-go/internal/domain/dtos/userdto"
	"github.com/tdatIT/backend-go/internal/tranport/http/helper"
	"github.com/tdatIT/backend-go/pkgs/utils/valid"
)

type AuthHandler struct {
	app *auth.Application
}

func NewAuthHandler(app *auth.Application) *AuthHandler {
	return &AuthHandler{app: app}
}

func (h *AuthHandler) LoginByUserPass(c *echo.Context) error {
	req := new(userdto.LoginByUserPassReq)
	if err := c.Bind(req); err != nil {
		slog.Error("failed to bind request body", slog.String("error", err.Error()))
		return err
	}

	if err := valid.GetValidator().Validate(req); err != nil {
		slog.Warn("validation request body failed", slog.String("error", err.Error()))
		return err
	}

	res, err := h.app.Queries.LoginByUsernameAndPassword.Handle(c.Request().Context(), req)
	if err != nil {
		slog.Error("failed to login by username and password", slog.String("error", err.Error()))
		return err
	}

	return helper.WriteSuccess(c, res)
}

func (h *AuthHandler) LoginByGoogle(c *echo.Context) error {
	req := new(userdto.LoginByGoogleReq)
	if err := c.Bind(req); err != nil {
		slog.Error("failed to bind request body", slog.String("error", err.Error()))
		return err
	}

	req.UserAgent = c.Request().UserAgent()
	req.IPAddress = c.RealIP()

	if err := valid.GetValidator().Validate(req); err != nil {
		slog.Warn("validation request body failed", slog.String("error", err.Error()))
		return err
	}

	res, err := h.app.Queries.LoginByGoogle.Handle(c.Request().Context(), req)
	if err != nil {
		slog.Error("failed to login by google", slog.String("error", err.Error()))
		return err
	}

	return helper.WriteSuccess(c, res)
}

func (h *AuthHandler) Register(c *echo.Context) error {
	req := new(userdto.RegisterReq)
	if err := c.Bind(req); err != nil {
		slog.Error("failed to bind request body", slog.String("error", err.Error()))
		return err
	}

	if err := valid.GetValidator().Validate(req); err != nil {
		slog.Warn("validation request body failed", slog.String("error", err.Error()))
		return err
	}

	res, err := h.app.Commands.Register.Handle(c.Request().Context(), req)
	if err != nil {
		slog.Error("failed to register", slog.String("error", err.Error()))
		return err
	}

	return helper.WriteSuccess(c, res)
}

func extractBearerToken(c *echo.Context) (string, error) {
	bearer := c.Request().Header.Get("Authorization")
	if bearer == "" {
		slog.Warn("missing Authorization header")
		return "", helper.ErrMissingAuthHeader
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(bearer, prefix) {
		slog.Warn("invalid Authorization header format")
		return "", helper.ErrInvalidAuthHeader
	}

	token := strings.TrimSpace(bearer[len(prefix):])
	if token == "" {
		slog.Warn("invalid Authorization header format")
		return "", helper.ErrInvalidAuthHeader
	}

	return token, nil
}

func (h *AuthHandler) RefreshToken(c *echo.Context) error {
	token, err := extractBearerToken(c)
	if err != nil {
		return err
	}

	req := &userdto.RefreshTokenReq{RefreshToken: token}

	res, err := h.app.Queries.RefreshToken.Handle(c.Request().Context(), req)
	if err != nil {
		slog.Error("failed to refresh token", slog.String("error", err.Error()))
		return err
	}

	return helper.WriteSuccess(c, res)
}

func (h *AuthHandler) Logout(c *echo.Context) error {
	token, err := extractBearerToken(c)
	if err != nil {
		return err
	}

	req := &userdto.LogoutReq{AccessToken: token}

	if err := h.app.Commands.Logout.Handle(c.Request().Context(), req); err != nil {
		slog.Error("failed to logout", slog.String("error", err.Error()))
		return err
	}

	return helper.WriteSuccess(c, nil)
}
