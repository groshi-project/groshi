package ghttp

import (
	"github.com/jieggii/groshi/internal/database"
)

type Handle func(request *Request, currentUser *database.User)
