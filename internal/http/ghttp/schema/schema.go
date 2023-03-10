package schema

// ErrorOrigin represents origin of error: server or client.
type ErrorOrigin string

const (
	ClientErrorOrigin ErrorOrigin = "client"
	ServerErrorOrigin ErrorOrigin = "server"
)

type ErrorTag string

const (
	// InvalidRequestErrorTag is returned when request did not pass validation.
	InvalidRequestErrorTag ErrorTag = "invalid_request"

	// UnauthorizedErrorTag is returned when request is not unauthorized (when it has to be)
	UnauthorizedErrorTag ErrorTag = "unauthorized"

	// InternalServerErrorErrorTag is returned when any internal server error happens.
	InternalServerErrorErrorTag ErrorTag = "internal_server_error"

	// AccessDeniedErrorTag is returned when user have no access to resource.
	AccessDeniedErrorTag ErrorTag = "access_denied"

	// ConflictErrorTag is returned when request causes any kind of conflict.
	ConflictErrorTag = "conflict"

	// ObjectNotFoundErrorTag is returned when object was not found.
	ObjectNotFoundErrorTag ErrorTag = "object_not_found"
)

type SuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

type ErrorResponse struct {
	Success bool `json:"success"`

	ErrorTag     ErrorTag    `json:"error_tag"`
	ErrorOrigin  ErrorOrigin `json:"error_origin"`
	ErrorDetails string      `json:"error_details"`
}
