package util

import (
	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/loggers"
	"net/http"
)

func abortWithErrorMessage(c *gin.Context, statusCode int, errorDetail string) {
	if statusCode == http.StatusInternalServerError {
		loggers.Error.Printf("internal server error: %v", errorDetail)
		errorDetail = "internal server error" // todo
	}
	c.AbortWithStatusJSON(statusCode, gin.H{
		"error_detail": errorDetail,
	})
}

func AbortWithStatusInternalServerError(c *gin.Context, err error) {
	abortWithErrorMessage(
		c, http.StatusInternalServerError, err.Error(),
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
