package security

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
)

type TokenManager interface {
	GenerateTokens(subjectID uint64, sessionID string, refreshJTI string) (accessToken string, refreshToken string, accessExp time.Time, err error)
	VerifyToken(tokenString string) (*CustomClaims, error)
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type JWTTokenManager struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type CustomClaims struct {
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

func NewJWTTokenManager(cfg JWTConfig) *JWTTokenManager {
	if cfg.Secret == "" {
		slog.Error("JWT secret is required")
		os.Exit(1)
	}

	accessTTL := cfg.AccessTokenTTL
	if accessTTL == 0 {
		accessTTL = 15 * time.Minute
	}

	refreshTTL := cfg.RefreshTokenTTL
	if refreshTTL == 0 {
		refreshTTL = 30 * 24 * time.Hour
	}

	return &JWTTokenManager{
		secret:          []byte(cfg.Secret),
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

func (m *JWTTokenManager) GenerateTokens(subjectID uint64, sessionID string, refreshJTI string) (string, string, time.Time, error) {
	now := time.Now()
	accessExp := now.Add(m.accessTokenTTL)
	refreshExp := now.Add(m.refreshTokenTTL)

	accessToken, err := m.signToken(subjectID, sessionID, uuid.NewString(), accessExp)
	if err != nil {
		return "", "", time.Time{}, err
	}

	refreshToken, err := m.signToken(subjectID, sessionID, refreshJTI, refreshExp)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken, refreshToken, accessExp, nil
}

func (m *JWTTokenManager) signToken(subjectID uint64, sessionID string, tokenID string, exp time.Time) (string, error) {
	claims := CustomClaims{
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(subjectID, 10),
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTTokenManager) VerifyToken(tokenString string) (*CustomClaims, error) {
	claims := new(CustomClaims)
	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
