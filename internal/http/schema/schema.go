package schema

// ErrorOrigin represents origin of error.
// It can only be server or client.
type ErrorOrigin string

const (
	ErrorOriginClient ErrorOrigin = "client"
	ErrorOriginServer ErrorOrigin = "server"
)

type SuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

type ErrorResponse struct {
	Success      bool        `json:"success"`
	ErrorMessage string      `json:"error_message"`
	ErrorOrigin  ErrorOrigin `json:"error_origin"`
}