// Package ghttp (stands for "groshi HTTP") is basically a tiny HTTP framework
// which simplifies reading and writing HTTP requests.
package ghttp

import (
	"bytes"
	"encoding/json"
	"github.com/jieggii/groshi/internal/http/schema"
	"github.com/jieggii/groshi/internal/loggers"
	"io"
	"net/http"
)

// Request is struct which keeps user ghttp data.
type Request struct {
	ResponseWriter http.ResponseWriter
	RawRequest     *http.Request
}

// NewRequest creates new Request instance.
func NewRequest(w http.ResponseWriter, r *http.Request) *Request {
	return &Request{ResponseWriter: w, RawRequest: r}
}

// sendJSONResponse sends JSON response.
func (req *Request) sendJSONResponse(data interface{}) {
	req.ResponseWriter.Header().Set("Content-Type", "application/json")
	req.ResponseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(req.ResponseWriter).Encode(data)
}

// Decode decodes ghttp body.
func (req *Request) Decode(v interface{}) error {
	body, err := io.ReadAll(req.RawRequest.Body)
	if err != nil {
		panic(err) // todo
	}
	// todo:
	req.RawRequest.Body = io.NopCloser(bytes.NewBuffer(body))
	err = json.NewDecoder(req.RawRequest.Body).Decode(&v)
	req.RawRequest.Body = io.NopCloser(bytes.NewBuffer(body))

	return err
}

// DecodeSafe decodes ghttp body and handles error.
// Returns true if there was no error.
// todo?: merge Decode and DecodeSafe
func (req *Request) DecodeSafe(v interface{}) bool {
	if err := req.Decode(v); err != nil {
		req.SendClientSideErrorResponse(schema.InvalidRequestBody)
		return false
	}
	return true
}

// SendSuccessResponse sends success response.
func (req *Request) SendSuccessResponse(data interface{}) {
	successObject := schema.SuccessResponse{Success: true, Data: data}
	req.sendJSONResponse(successObject)
}

// SendClientSideErrorResponse sends client-side error response containing error message.
func (req *Request) SendClientSideErrorResponse(errorMessage string) {
	req.sendJSONResponse(schema.ErrorResponse{
		Success:      false,
		ErrorOrigin:  schema.ErrorOriginClient,
		ErrorMessage: errorMessage,
	})
}

// SendServerSideErrorResponse sends server-side error response containing error message
// "Internal Server Error", logs error comment and error object.
func (req *Request) SendServerSideErrorResponse(errorComment string, err error) {
	req.sendJSONResponse(schema.ErrorResponse{
		Success:      false,
		ErrorOrigin:  schema.ErrorOriginServer,
		ErrorMessage: schema.InternalServerError,
	})
	loggers.Error.Printf("%v (%v)", errorComment, err)
}
