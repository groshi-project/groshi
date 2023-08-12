package util

import (
	"github.com/gin-gonic/gin"
	"github.com/jieggii/groshi/internal/loggers"
)

// BindBody is an alias function for gin.Context.ShouldBind to be used inside handlers.
func BindBody(c *gin.Context, v interface{}) (ok bool) {
	if err := c.ShouldBind(v); err != nil {
		AbortWithStatusBadRequest(c, "invalid request body params")
		return false
	}
	return true
}

// BindQuery is an alias for gin.Context.ShouldBindQuery to be used inside handlers.
func BindQuery(c *gin.Context, v interface{}) (ok bool) {
	if err := c.ShouldBindQuery(v); err != nil {
		loggers.Error.Printf("error binding query: %v", err)
		AbortWithStatusBadRequest(c, "invalid query params")
		return false
	}
	return true
}
