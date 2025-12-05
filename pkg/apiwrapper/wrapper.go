package apiwrapper

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/leehai1107/cmm_server/pkg/errors"
)

const (
	statusSuccess = 1
	statusFail    = 0
)

/*
*
Example success response:

	{
	    "status": 1,
	    "code": 1,
	    "message": "Thành công!",
		"error": null,
	    "data": {
	        "abc": {
	            "cde": "ikh",
	            "fgh": "jkl"
	        }
	    },
	}

Example failure response:

	{
	    "status": 0,
	    "code": -1,
	    "message": "Yêu cầu không hợp lệ!",
	    "error": {
	        "code": -1,
	        "message": "Yêu cầu không hợp lệ!"
	    },
		"result":{
			"data":{}
		}
	}
*/
type APIResponse struct {
	Status  int64            `json:"status"`
	Code    errors.ErrorType `json:"code"`
	Message string           `json:"message"`
	Data    interface{}      `json:"data,omitempty"`
	Error   interface{}      `json:"error,omitempty"`
}
type apiHandlerFn func(c *gin.Context) *APIResponse

func Wrap(fn apiHandlerFn) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiRsp := fn(c)
		if apiRsp != nil {
			if apiRsp.Status == statusSuccess {
				c.JSON(http.StatusOK, apiRsp)
				return
			}
			// Use appropriate status code based on error type
			statusCode := getStatusCodeFromErrorType(apiRsp.Code)
			c.JSON(statusCode, apiRsp)
		}
	}
}

func Abort(c *gin.Context, apiRsp *APIResponse) {
	if apiRsp.Status == statusSuccess {
		c.AbortWithStatusJSON(http.StatusOK, apiRsp)
		return
	}
	statusCode := getStatusCodeFromErrorType(apiRsp.Code)
	c.AbortWithStatusJSON(statusCode, apiRsp)
}

func AbortWithStatus(c *gin.Context, statusCode int, apiRsp *APIResponse) {
	c.AbortWithStatusJSON(statusCode, apiRsp)
}

func File(fn apiHandlerFn) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiRsp := fn(c)
		if apiRsp.Status == statusSuccess {
			fileName := apiRsp.Data.(string)
			c.FileAttachment(fileName, fileName)
			os.Remove(fileName)
			return
		}
		statusCode := getStatusCodeFromErrorType(apiRsp.Code)
		c.JSON(statusCode, "error")
	}
}

// Helper function to determine HTTP status code from error type
func getStatusCodeFromErrorType(errType errors.ErrorType) int {
	switch errType {
	case errors.BadRequestErr:
		return http.StatusBadRequest
	case errors.NotFound:
		return http.StatusNotFound
	case errors.AuthenticationFailed:
		return http.StatusUnauthorized
	case errors.InternalServerError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
