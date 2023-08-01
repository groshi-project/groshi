package utils

import (
	"github.com/gin-gonic/gin"
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
