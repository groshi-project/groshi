package ghttp

import "github.com/jieggii/groshi/internal/database"

// Handle is type of handle which is used to define all groshi handles.
type Handle func(request *Request, currentUser *database.User)

// _emptyStruct is well... An empty struct!
type _emptyStruct struct{}

// EmptyRequest is type used to define requests without parameters.
type EmptyRequest = _emptyStruct

// EmptyResponse is type used to define responses without any data.
type EmptyResponse = _emptyStruct
