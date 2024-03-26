package httpresp

import (
	"encoding/json"
	"net/http"
)

// Response represents HTTP response.
type Response struct {
	// HTTP response status code.
	statusCode int

	// Response body.
	body any
}

// New creates a new instance of [Response] and returns pointer to it.
func New(statusCode int, body any) *Response {
	return &Response{statusCode: statusCode, body: body}
}

// NewOK creates a new instance of [Response] with HTTP status code 200 and returns pointer to it.
func NewOK(body any) *Response {
	return &Response{statusCode: http.StatusOK, body: body}
}

// Render renders HTTP response.
// renderJSON sends response header with provided statusCode,
// sets Content-Type header key to "application/json" and writes json encoding of v to the provided [http.ResponseWriter] w.
func renderJSON(w http.ResponseWriter, statusCode int, v any) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		panic(err)
	}
}

// Render renders HTTP response.
func Render(w http.ResponseWriter, resp *Response) {
	renderJSON(w, resp.statusCode, resp.body)
}
