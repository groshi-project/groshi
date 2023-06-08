package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/loggers"
	"net/http"
)

// SendErrorResponse must be used for any error to be sent to user except internal server error.
// (SendInternalServerError must be used for internal server error).
func SendErrorResponse(c *gin.Context, statusCode int, errorMessage string) {
	c.JSON(statusCode, gin.H{"error_message": errorMessage})
}

// SendInternalServerErrorResponse must be used to send internal server error to the client
// and in order to log internal server error.
func SendInternalServerErrorResponse(c *gin.Context, logMessage string, err error) {
	loggers.Error.Printf("Returned internal server error: %v (%v).", logMessage, err)
	c.JSON(http.StatusInternalServerError, gin.H{"error_message": "internal server error"})
}
