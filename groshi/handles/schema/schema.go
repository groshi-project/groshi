package schema

type SuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

type ErrorCode int

var (
	ClientSideError ErrorCode = 1
	ServerSideError ErrorCode = 2
)

type ErrorResponse struct {
	Success      bool      `json:"success"`
	ErrorMessage string    `json:"error_message"`
	ErrorCode    ErrorCode `json:"error_code"`
}
