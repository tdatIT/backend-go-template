package helper

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/tdatIT/backend-go/pkgs/svcerr"
)

type Response struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	VIMessage string `json:"vi_message"`
	Data      any    `json:"data,omitempty"`
}

func WriteSuccess(c *echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Code:      "00",
		Message:   "success",
		VIMessage: "Thành công",
		Data:      data,
	})
}

func WriteError(c *echo.Context, err error) error {
	if svcErr, ok := errors.AsType[*svcerr.Error](err); ok {
		status := svcErr.HTTPStatus
		if status == 0 {
			status = http.StatusInternalServerError
		}

		return c.JSON(status, Response{
			Code:      svcErr.Code,
			Message:   svcErr.Message,
			VIMessage: svcErr.VIMessage,
			Data:      nil,
		})
	}

	if _, ok := errors.AsType[validator.ValidationErrors](err); ok {
		return c.JSON(http.StatusBadRequest, Response{
			Code:      "01",
			Message:   "invalid data",
			VIMessage: "Dữ liệu không hợp lệ",
			Data:      nil,
		})
	}

	var sc echo.HTTPStatusCoder
	if errors.As(err, &sc) { // find error in an error chain that implements HTTPStatusCoder
		if tmp := sc.StatusCode(); tmp != 0 {
			return c.JSON(tmp, Response{
				Code:      "01",
				Message:   err.Error(),
				VIMessage: "Lỗi xử lý",
				Data:      nil,
			})
		}
	}

	return c.JSON(http.StatusInternalServerError, Response{
		Code:      "01",
		Message:   "internal server error",
		VIMessage: "Lỗi máy chủ nội bộ",
		Data:      nil,
	})
}
