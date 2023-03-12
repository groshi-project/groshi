package ghttp

import "github.com/jieggii/groshi/internal/database"

// Handle is type of handle which is used to define all groshi handles.
type Handle func(request *Request, currentUser *database.User)

// RequestParams is interface for defining HTTP request parameters.
type RequestParams interface {
	Validate() bool
}

// EmptyRequestParams is type used to define requests without parameters.
type EmptyRequestParams = struct{}

func (p *EmptyRequestParams) Validate() bool {
	return true
}

// Response is interface for defining HTTP responses.
type Response interface{}

// EmptyResponse is, well... an empty response!
type EmptyResponse = struct{}
