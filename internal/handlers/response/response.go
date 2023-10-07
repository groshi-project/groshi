package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/groshi-project/groshi/internal/loggers"
	"github.com/groshi-project/groshi/internal/models"
)

// EmptyErrorDetails is used when aborting with error
// without error details not to create multiple empty slices.
var EmptyErrorDetails = make([]string, 0)

func abortWithErrorMessage(c *gin.Context, statusCode int, errorMessage string, errorDetails []string) {
	c.AbortWithStatusJSON(statusCode, models.Error{
		ErrorMessage: errorMessage,
		ErrorDetails: errorDetails,
	})
}

func AbortWithStatusInternalServerError(c *gin.Context, err error) {
	loggers.Error.Printf("aborted with internal server error: %v", err)
	abortWithErrorMessage(
		c, http.StatusInternalServerError, "internal server error", EmptyErrorDetails,
	)
}

func AbortWithStatusNotFound(c *gin.Context, errorMessage string) {
	abortWithErrorMessage(
		c, http.StatusNotFound, errorMessage, EmptyErrorDetails,
	)
}

func AbortWithStatusForbidden(c *gin.Context, errorMessage string) {
	abortWithErrorMessage(
		c, http.StatusForbidden, errorMessage, EmptyErrorDetails,
	)
}

func AbortWithStatusBadRequest(c *gin.Context, errorMessage string, errorDetails []string) {
	abortWithErrorMessage(
		c, http.StatusBadRequest, errorMessage, errorDetails,
	)
}

func AbortWithStatusConflict(c *gin.Context, errorMessage string) {
	abortWithErrorMessage(
		c, http.StatusConflict, errorMessage, EmptyErrorDetails,
	)
}

func AbortWithStatusUnauthorized(c *gin.Context, errorMessage string) {
	abortWithErrorMessage(c, http.StatusUnauthorized, errorMessage, EmptyErrorDetails)
}

func ReturnSuccessfulResponse(c *gin.Context, response interface{}) {
	c.JSON(http.StatusOK, response)
}
