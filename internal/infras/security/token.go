package security

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager interface {
	GenerateTokens(subjectID uint64) (accessToken string, refreshToken string, accessExp time.Time, err error)
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
	SubID uint64 `json:"sub_id"`
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

func (m *JWTTokenManager) GenerateTokens(subjectID uint64) (string, string, time.Time, error) {
	now := time.Now()
	accessExp := now.Add(m.accessTokenTTL)
	refreshExp := now.Add(m.refreshTokenTTL)

	accessToken, err := m.signToken(subjectID, accessExp)
	if err != nil {
		return "", "", time.Time{}, err
	}

	refreshToken, err := m.signToken(subjectID, refreshExp)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken, refreshToken, accessExp, nil
}

func (m *JWTTokenManager) signToken(subjectID uint64, exp time.Time) (string, error) {
	claims := CustomClaims{
		SubID: subjectID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
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
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
