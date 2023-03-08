package ghttp

import "github.com/jieggii/groshi/internal/database"

// Handle is type of handle which is used to define all groshi handles.
type Handle func(request *Request, currentUser *database.User)
