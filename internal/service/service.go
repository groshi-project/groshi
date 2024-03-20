package service

import (
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/service/handler"
	"github.com/groshi-project/groshi/internal/service/job"
	"github.com/groshi-project/groshi/pkg/jwtauthority"
	"github.com/groshi-project/groshi/pkg/passwdauthority"
	"log"
)

// Service represents groshi service containing all its dependencies.
type Service struct {
	// Service's HTTP handlers and their dependencies.
	Handler *handler.Handler

	// Enable Swagger UI route.
	Swagger bool

	// Service's jobs and their dependencies.
	job *job.Job
}

// New creates a new instance of [Service] and returns pointer to it.
func New(database *database.Database, jwtAuthority *jwtauthority.DefaultAuthority, passwordAuthority *passwdauthority.Authority, internalServerErrorLogger *log.Logger, swagger bool) *Service {
	return &Service{
		Handler: handler.New(database, jwtAuthority, passwordAuthority, internalServerErrorLogger),
		Swagger: swagger,
		job:     job.New(database),
	}
}
