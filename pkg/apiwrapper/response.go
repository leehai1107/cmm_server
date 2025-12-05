package apiwrapper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leehai1107/cmm_server/pkg/errors"
)

// Direct APIResponse helper functions
func SuccessAPIResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Status:  statusSuccess,
		Code:    errors.Success,
		Message: errors.GetMessage(errors.Success.New()),
		Data:    data,
		Error:   nil,
	}
}

func ErrorAPIResponse(code errors.ErrorType, message string) *APIResponse {
	return &APIResponse{
		Status:  statusFail,
		Code:    code,
		Message: message,
		Error: map[string]interface{}{
			"code":    code,
			"message": message,
		},
		Data: nil,
	}
}

// Helper functions to send API responses directly
func SendSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessAPIResponse(data))
}

func SendError(c *gin.Context, statusCode int, errType errors.ErrorType, message string) {
	c.AbortWithStatusJSON(statusCode, ErrorAPIResponse(errType, message))
}

// Common HTTP status code helpers
func SendBadRequest(c *gin.Context, message string) {
	SendError(c, http.StatusBadRequest, errors.BadRequestErr, message)
}

func SendUnauthorized(c *gin.Context, message string) {
	SendError(c, http.StatusUnauthorized, errors.AuthenticationFailed, message)
}

func SendInternalError(c *gin.Context, message string) {
	SendError(c, http.StatusInternalServerError, errors.InternalServerError, message)
}
