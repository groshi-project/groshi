package service

import (
	"github.com/groshi-project/groshi/internal/auth"
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/service/handler"
	"github.com/groshi-project/groshi/internal/service/job"
	"log"
)

// Service represents groshi service containing all its dependencies.
type Service struct {
	// HTTP handlers and their dependencies.
	Handler *handler.Handler

	// Enable Swagger UI route.
	SwaggerEnable bool

	// Jobs and their dependencies.
	job *job.Job
}

// New creates a new instance of [Service] and returns pointer to it.
func New(database database.Database, jwtAuthenticator auth.JWTAuthenticator, passwordAuthenticator auth.PasswordAuthenticator, internalServerErrorLogger *log.Logger, swagger bool) *Service {
	return &Service{
		Handler:       handler.New(database, jwtAuthenticator, passwordAuthenticator, internalServerErrorLogger),
		SwaggerEnable: swagger,
		job:           job.New(database),
	}
}
