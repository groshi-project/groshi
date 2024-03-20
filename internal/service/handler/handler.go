package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/pkg/jwtauthority"
	"github.com/groshi-project/groshi/pkg/passwdauthority"
	"log"
)

// Handler represents dependencies for HTTP handler functions.
type Handler struct {
	// database used to store and retrieve data.
	database *database.Database

	// JWTAuthority is a JWT authority used to generate and validate JavaScript Web Tokens.
	JWTAuthority *jwtauthority.DefaultAuthority

	// passwordAuthority is a password authority used to hash and validate passwords.
	passwordAuthority *passwdauthority.Authority

	// internalServerErrorLogger is an internal server error logger used to log internal server errors :).
	internalServerErrorLogger *log.Logger

	// paramsValidate contains validator settings for validating incoming request params.
	paramsValidate *validator.Validate
}

// New creates a new instance of [Handler] and returns pointer to it.
func New(database *database.Database, jwtAuthority *jwtauthority.DefaultAuthority, passwordAuthority *passwdauthority.Authority, internalServerErrorLogger *log.Logger) *Handler {
	return &Handler{
		database:                  database,
		JWTAuthority:              jwtAuthority,
		passwordAuthority:         passwordAuthority,
		internalServerErrorLogger: internalServerErrorLogger,
		paramsValidate:            validator.New(),
	}
}
