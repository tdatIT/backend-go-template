package query

import (
	"context"

	"github.com/tdatIT/backend-go/internal/application/auth/helper"
	"github.com/tdatIT/backend-go/internal/domain/dtos/userdto"
	"github.com/tdatIT/backend-go/internal/infras/security"
	"github.com/tdatIT/backend-go/pkgs/decorator"
)

type IVerifyTokenQuery decorator.QueryHandler[*userdto.VerifyTokenReq, *userdto.VerifyTokenRes]

type verifyTokenQuery struct {
	tokenManager security.TokenManager
}

func NewVerifyTokenQuery(tokenManager security.TokenManager) IVerifyTokenQuery {
	return &verifyTokenQuery{tokenManager: tokenManager}
}

func (v verifyTokenQuery) Handle(ctx context.Context, req *userdto.VerifyTokenReq) (*userdto.VerifyTokenRes, error) {
	claims, err := v.tokenManager.VerifyToken(req.AccessToken)
	if err != nil {
		return nil, helper.ErrInvalidToken
	}

	return &userdto.VerifyTokenRes{
		SubID:     claims.SubID,
		JTI:       claims.ID,
		ExpiresAt: claims.ExpiresAt.Time.Unix(),
	}, nil
}
