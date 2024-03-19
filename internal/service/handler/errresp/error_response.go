package errresp

// ErrorResponse represents error response model.
type ErrorResponse struct {
	ErrorMessage string `json:"error_message" example:"today is a sunny day so I decided to go for a walk instead of serving your requests"`
}

// NewErrorResponse creates a new instance of [ErrorResponse] and returns pointer to it.
func NewErrorResponse(errorMessage string) *ErrorResponse {
	return &ErrorResponse{ErrorMessage: errorMessage}
}
