package util

import (
	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/loggers"
	"net/http"
)

// emptyErrorDetails is used when aborting with error
// without error details not to create multiple empty slices.
var emptyErrorDetails = make([]string, 0)

func abortWithErrorMessage(c *gin.Context, statusCode int, errorDescription string, errorDetails []string) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"error_description": errorDescription,
		"error_details":     errorDetails,
	})
}

func AbortWithStatusInternalServerError(c *gin.Context, err error) {
	loggers.Error.Printf("aborted with internal server error: %v", err)
	abortWithErrorMessage(
		c, http.StatusInternalServerError, "internal server error", emptyErrorDetails,
	)
}

func AbortWithStatusNotFound(c *gin.Context, errorDescription string) {
	abortWithErrorMessage(
		c, http.StatusNotFound, errorDescription, emptyErrorDetails,
	)
}

func AbortWithStatusForbidden(c *gin.Context, errorDescription string) {
	abortWithErrorMessage(
		c, http.StatusForbidden, errorDescription, emptyErrorDetails,
	)
}

func AbortWithStatusBadRequest(c *gin.Context, errorDescription string, errorDetails []string) {
	abortWithErrorMessage(
		c, http.StatusBadRequest, errorDescription, errorDetails,
	)
}

func AbortWithStatusConflict(c *gin.Context, errorDescription string) {
	abortWithErrorMessage(
		c, http.StatusConflict, errorDescription, emptyErrorDetails,
	)
}

func ReturnSuccessfulResponse(c *gin.Context, response interface{}) {
	c.JSON(http.StatusOK, response)
}
