package ghttp

import (
	"github.com/jieggii/groshi/groshi/handles/jwt"
)

type Handle func(req *Request, claims *jwt.Claims)

//type AuthorizedHandle func(req *Request, user database.User, claims *jwt.Claims)
