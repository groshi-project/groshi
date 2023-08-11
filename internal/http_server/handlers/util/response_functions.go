package util

import (
	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/loggers"
	"net/http"
)

// AbortWithErrorMessage TODO
func AbortWithErrorMessage(c *gin.Context, statusCode int, errorDetail string) {
	if statusCode == http.StatusInternalServerError {
		loggers.Error.Printf("internal server error: %v", errorDetail)
		errorDetail = "internal server error" // todo
	}
	c.AbortWithStatusJSON(statusCode, gin.H{
		"error_detail": errorDetail,
	})
}

// AbortWithInternalServerError aborts with internal server error.
func AbortWithInternalServerError(c *gin.Context, err error) {
	AbortWithErrorMessage(
		c, http.StatusInternalServerError, err.Error(),
	)
}

func AbortWithNotFoundError(c *gin.Context, errorDetail string) {
	AbortWithErrorMessage(
		c, http.StatusNotFound, errorDetail,
	)
}

func AbortWithForbiddenError(c *gin.Context, errorDetail string) {
	AbortWithErrorMessage(
		c, http.StatusForbidden, errorDetail,
	)
}

func AbortWithBadRequest(c *gin.Context, errorDetail string) {
	AbortWithErrorMessage(
		c, http.StatusBadRequest, errorDetail, // todo
	)
}

func AbortWithConflictError(c *gin.Context, errorDetail string) {
	AbortWithErrorMessage(
		c, http.StatusConflict, errorDetail,
	)
}

func ReturnSuccessfulResponse(c *gin.Context, response interface{}) {
	c.JSON(http.StatusOK, response)
}
