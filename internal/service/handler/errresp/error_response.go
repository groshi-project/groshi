package errresp

// ErrorData represents error response model.
type ErrorData struct {
	ErrorMessage string `json:"error_message"`
}

// NewErrorData creates a new instance of [ErrorData] and returns pointer to it.
func NewErrorData(errorMessage string) *ErrorData {
	return &ErrorData{ErrorMessage: errorMessage}
}
