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
	// handler contains groshi's HTTP handlers and their dependencies.
	handler *handler.Handler

	// job contains groshi's jobs and their dependencies.
	job *job.Job
}

// New creates a new instance of [Service] and returns pointer to it.
func New(database *database.Database, jwtAuthority *jwtauthority.JWTAuthority, passwordAuthority *passwdauthority.PasswordAuthority, internalServerErrorLogger *log.Logger) *Service {
	return &Service{
		handler: handler.New(database, jwtAuthority, passwordAuthority, internalServerErrorLogger),
		job:     job.New(database),
	}
}
