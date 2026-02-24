package helper

import (
	"net/http"

	"github.com/tdatIT/backend-go/pkgs/svcerr"
	"google.golang.org/grpc/codes"
)

var (
	ErrInvalidUserOrPwd = &svcerr.SvcErr{
		Message:    "Invalid username or password",
		VIMessage:  "Tên đăng nhập hoặc mật khẩu không hợp lệ",
		Code:       "AUTH-001",
		HTTPStatus: http.StatusUnauthorized,
		GRPCCode:   codes.Unauthenticated,
	}

	ErrInvalidToken = &svcerr.SvcErr{
		Message:    "Invalid token",
		VIMessage:  "Token không hợp lệ",
		Code:       "AUTH-002",
		HTTPStatus: http.StatusUnauthorized,
		GRPCCode:   codes.Unauthenticated,
	}

	ErrUserAlreadyExists = &svcerr.SvcErr{
		Message:    "User already exists",
		VIMessage:  "Nguoi dung da ton tai",
		Code:       "AUTH-003",
		HTTPStatus: http.StatusConflict,
		GRPCCode:   codes.AlreadyExists,
	}
)
