package util

import (
	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/http_server/error_messages"
	"github.com/jieggii/groshi/internal/loggers"
	"net/http"
)

// AbortWithErrorMessage aborts with provided status code and returns JSON response containing
// error message.
// Logs error and returns status code 500 with the default error message ("internal server error")
// if status code is 500.
func AbortWithErrorMessage(c *gin.Context, statusCode int, errorMessage string) {
	if statusCode == http.StatusInternalServerError {
		loggers.Error.Printf("internal server error: %v", errorMessage)
		errorMessage = "internal server error"
	}
	c.AbortWithStatusJSON(statusCode, gin.H{
		"error_message": errorMessage,
	})
}

// AbortWithInternalServerError aborts with internal server error.
func AbortWithInternalServerError(c *gin.Context, err error) {
	AbortWithErrorMessage(c, http.StatusInternalServerError, err.Error())
}

// ReturnSuccessfulResponse TODO
func ReturnSuccessfulResponse(c *gin.Context, response interface{}) {
	c.JSON(http.StatusOK, response)
}

// BindBody is an alias function for gin.Context.ShouldBind to be used inside handlers.
func BindBody(c *gin.Context, v interface{}) (ok bool) {
	if err := c.ShouldBind(v); err != nil {
		AbortWithErrorMessage(
			c, http.StatusBadRequest, error_messages.ErrorInvalidRequestParams.Error(),
		)
		return false
	}
	return true
}

func BindQuery(c *gin.Context, v interface{}) (ok bool) {
	if err := c.ShouldBindQuery(v); err != nil {
		AbortWithErrorMessage(
			c, http.StatusBadRequest, error_messages.ErrorInvalidRequestParams.Error(),
		)
		return false
	}
	return true
}
