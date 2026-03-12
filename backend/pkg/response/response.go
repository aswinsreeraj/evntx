package response

import (
	"errors"
	"net/http"

	pkgerrors "github.com/aswinsreeraj/evntx/pkg/errors"
	"github.com/gin-gonic/gin"
)

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type APIResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	Data    interface{}  `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(200, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Error:   nil,
	})
}

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

// AppError safely unwraps and responds with a custom AppError, or defaults to 500 if unknown.
func AppError(c *gin.Context, err error) {
	var appErr *pkgerrors.AppError
	if errors.As(err, &appErr) {
		Error(c, appErr.HTTPCode, appErr.Code, appErr.Message)
		return
	}

	// Fallback for unknown errors
	Error(c, http.StatusInternalServerError, pkgerrors.InternalServerError, "Internal server error")
}
