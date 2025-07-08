package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ErrCodeInvalidRequest = "INVALID_REQUEST"
	ErrCodeInternalServer = "INTERNAL_SERVER_ERROR"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeValidation     = "VALIDATION_ERROR"
)

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

func BadRequest(c *gin.Context, code, message, details string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func ValidationError(c *gin.Context, details string) {
	BadRequest(c, ErrCodeValidation, "Request parameters are invalid", details)
}
