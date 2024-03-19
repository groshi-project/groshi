package model

// Error represents error response model.
type Error struct {
	ErrorMessage string `json:"error_message" example:"today is a sunny day so I decided to go for a walk instead of serving your requests"`
}

// NewError creates a new instance of [Error] and returns pointer to it.
func NewError(errorMessage string) *Error {
	return &Error{ErrorMessage: errorMessage}
}
