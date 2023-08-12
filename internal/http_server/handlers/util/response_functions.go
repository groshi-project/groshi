package util

import (
	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/loggers"
	"net/http"
)

func abortWithErrorMessage(c *gin.Context, statusCode int, errorDetail string) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"error_detail": errorDetail,
	})
}

func AbortWithStatusInternalServerError(c *gin.Context, err error) {
	loggers.Error.Printf("aborted with internal server error: %v", err)
	abortWithErrorMessage(
		c, http.StatusInternalServerError, "internal server error",
	)
}

func AbortWithStatusNotFound(c *gin.Context, errorDetail string) {
	abortWithErrorMessage(
		c, http.StatusNotFound, errorDetail,
	)
}

func AbortWithStatusForbidden(c *gin.Context, errorDetail string) {
	abortWithErrorMessage(
		c, http.StatusForbidden, errorDetail,
	)
}

func AbortWithStatusBadRequest(c *gin.Context, errorDetail string) {
	abortWithErrorMessage(
		c, http.StatusBadRequest, errorDetail,
	)
}

func AbortWithStatusConflict(c *gin.Context, errorDetail string) {
	abortWithErrorMessage(
		c, http.StatusConflict, errorDetail,
	)
}

func ReturnSuccessfulResponse(c *gin.Context, response interface{}) {
	c.JSON(http.StatusOK, response)
}
