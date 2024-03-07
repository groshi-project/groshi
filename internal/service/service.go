package service

import (
	"github.com/groshi-project/groshi/internal/database"
	"github.com/groshi-project/groshi/internal/service/authorities/jwt"
	"github.com/groshi-project/groshi/internal/service/authorities/passwd"
)

// Service represents groshi service.
type Service struct {
	Database          *database.Database
	PasswordAuthority *passwd.Authority
	JWTAuthority      *jwt.Authority
}

//func New(database *database.Database) *Service {
//	return &Service{Database: database}
//}
