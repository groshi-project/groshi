package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/groshi-project/groshi/internal/auth"
	"github.com/groshi-project/groshi/internal/database"
	"log"
)

// Handler represents dependencies for HTTP handler functions.
type Handler struct {
	// Database used to store and retrieve data.
	database *database.Database

	// JWT authenticator used to generate and validate JWTs.
	JWTAuthenticator auth.JWTAuthenticator

	// Password authenticator used to hash and validate passwords.
	passwordAuthenticator *auth.PasswordAuthenticator

	// Logger used to log internal server errors.
	internalServerErrorLogger *log.Logger

	// Validator settings for validating incoming request params.
	paramsValidate *validator.Validate
}

// New creates a new instance of [Handler] and returns pointer to it.
func New(database *database.Database, jwtAuthenticator auth.JWTAuthenticator, passwordAuthenticator *auth.PasswordAuthenticator, internalServerErrorLogger *log.Logger) *Handler {
	return &Handler{
		database:                  database,
		JWTAuthenticator:          jwtAuthenticator,
		passwordAuthenticator:     passwordAuthenticator,
		internalServerErrorLogger: internalServerErrorLogger,
		paramsValidate:            validator.New(),
	}
}
