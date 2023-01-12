package ghttp

import (
	"github.com/jieggii/groshi/groshi/database"
)

type Handle func(request *Request, currentUser *database.User)
