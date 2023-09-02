package util

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/groshi-project/groshi/internal/loggers"
)

func generateErrorDetails(err error) []string {
	errorDetails := make([]string, 0)

	if errors.As(err, new(validator.ValidationErrors)) { // if error is validator.ValidationErrors
		for _, fieldErr := range err.(validator.ValidationErrors) {
			errorDetails = append(errorDetails, fieldErr.Error())
		}
	} else { // if it is not! For example, *time.ParseError
		errorDetails = append(errorDetails, err.Error())
	}

	return errorDetails
}

// BindBody is an alias function for gin.Context.ShouldBind to be used inside handlers.
func BindBody(c *gin.Context, v interface{}) (ok bool) {
	if err := c.ShouldBind(v); err != nil {
		loggers.Error.Printf("error binding body: %v", err)

		AbortWithStatusBadRequest(
			c,
			"invalid request body params, please refer to the method documentation",
			generateErrorDetails(err),
		)
		return false
	}
	return true
}

// BindQuery is an alias for gin.Context.ShouldBindQuery to be used inside handlers.
func BindQuery(c *gin.Context, v interface{}) (ok bool) {
	if err := c.ShouldBindQuery(v); err != nil {
		loggers.Error.Printf("error binding query: %v", err)

		AbortWithStatusBadRequest(
			c,
			"invalid query params, please refer to the method documentation",
			generateErrorDetails(err),
		)
		return false
	}
	return true
}
