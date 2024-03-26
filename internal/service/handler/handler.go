package handler

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/groshi-project/groshi/internal/auth"
	"github.com/groshi-project/groshi/internal/database"
	"log"
)

var errMissingUsernameContextValue = errors.New("missing username context value")

// Handler represents dependencies for HTTP handler functions.
type Handler struct {
	// DefaultDatabase used to store and retrieve data.
	database database.Database

	// JWT authenticator used to generate and validate JWTs.
	JWTAuth auth.JWTAuthenticator

	// Password authenticator used to hash and validate passwords.
	passwordAuth auth.PasswordAuthenticator

	// Logger used to log internal server errors.
	internalServerErrorLogger *log.Logger

	// Validator settings for validating incoming request params.
	paramsValidate *validator.Validate
}

// New creates a new instance of [Handler] and returns pointer to it.
func New(database database.Database, jwtAuth auth.JWTAuthenticator, passwordAuth auth.PasswordAuthenticator, internalServerErrorLogger *log.Logger) *Handler {
	return &Handler{
		database:                  database,
		JWTAuth:                   jwtAuth,
		passwordAuth:              passwordAuth,
		internalServerErrorLogger: internalServerErrorLogger,
		paramsValidate:            validator.New(),
	}
}
