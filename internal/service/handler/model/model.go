package model

// Error represents error response model.
type Error struct {
	ErrorMessage string `json:"error_message" example:"example error message (who cares)"`
}

// NewError creates a new instance of [Error] and returns pointer to it.
func NewError(errorMessage string) *Error {
	return &Error{ErrorMessage: errorMessage}
}
