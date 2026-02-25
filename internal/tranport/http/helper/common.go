package helper

import (
	"net/http"

	"github.com/tdatIT/backend-go/pkgs/svcerr"
)

var (
	ErrMissingAuthHeader = &svcerr.Error{
		Message:    "missing authorization header",
		VIMessage:  "Thiếu header xác thực",
		Code:       "01",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrInvalidAuthHeader = &svcerr.Error{
		Message:    "invalid authorization header format",
		VIMessage:  "Định dạng header xác thực không hợp lệ",
		Code:       "01",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrRequestTimeout = &svcerr.Error{
		Message:    "request timeout",
		VIMessage:  "Yêu cầu đã hết thời gian chờ",
		Code:       "01",
		HTTPStatus: http.StatusGatewayTimeout,
	}
)
